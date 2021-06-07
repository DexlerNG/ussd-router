package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"time"
	"ussd-router/controllers"
	"ussd-router/controllers/configuration"
	"ussd-router/controllers/receive"
	"ussd-router/controllers/send"
	"ussd-router/startups"
	redis "ussd-router/startups/cache"
	"ussd-router/startups/queues"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	error := c.JSON(code, map[string]interface{}{
		"error": err.Error(),
	})
	if error != nil {
		c.Logger().Error(error)
	}
	//c.Logger().Error(err)
}

func authenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var startTime = time.Now()

		//log start time
		//fmt.Println("Authenticating")
		//if !utils.IsStringEmpty(c.QueryParam("clientId")) {
		//	c.Set("clientId", c.QueryParam("clientId"))
		//} else {
		//	c.Set("clientId", c.Request().Header.Get("client-id"))
		//}

		fmt.Println("Calling Function")
		if err := next(c); err != nil{
			c.Error(err)
		}
		fmt.Println("After Calling Function")

		fmt.Println("Headers", c.Request().Header)
		//fmt.Println("Request Body", string(reqBody))
		//fmt.Println("Response Body", string(resBody))
		fmt.Println("Host", c.Request().Host)
		fmt.Println("Method", c.Request().Method)
		fmt.Println("IP", c.RealIP())
		fmt.Println("Status Code", c.Response().Status)
		fmt.Println("URL", c.Request().URL)
		fmt.Println("Response Time Spent", time.Since(startTime))
		//log end time
		return nil
	}
}




func main() {
	err := godotenv.Load()
	fmt.Println("Godotenv Error:", err)
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler

	client := startups.GetMongoClient()
	defer client.Disconnect(context.Background())

	redis.GetRedisClient()
	queues.GetRabbitMQClient()
	defer queues.GetRabbitMQConnection().Close()


	requestLoggerQueueName := "elasticsearch.single.runner"
	requestLoggerIndexName := "ussd.router.request.logs." + os.Getenv("APP_ENV")
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		//fmt.Println("p", c.Request().Response.Header)
		go func() {
			logData := map[string]interface{}{
				"host":               c.Request().Host,
				"headers":            c.Request().Header,
				"method":             c.Request().Method,
				"ip":                 c.RealIP(),
				"url":                c.Request().URL.Path,
				"urlHost":            c.Request().URL.Host,
				"urlHostName":        c.Request().URL.Hostname(),
				"requestBody":        string(reqBody),
				"responseBody":       string(resBody),
				"responseType":       c.Response().Header().Get("Content-Type"),
				"responseStatusCode": c.Response().Status,
				"time": time.Now().Unix(),
			}

			//jsonValue, _ := json.Marshal(logData)
			//fmt.Println("Request Body", string(jsonValue))
			err = queues.RabbitMQPublishToQueue(requestLoggerQueueName, map[string]interface{}{
				"index": requestLoggerIndexName,
				"data":  logData,
			})
			if err != nil {
				fmt.Println("logged write error", err)
			}
		}()
	}))
	e.Any("/health-check", controllers.HealthCheckHandler)

	//e.Use()

	authorisedRoute := e.Group("/v1")
	//authorisedRoute.Use(authenticateMiddleware)
	//e.Use(authenticateMiddleware)

	authorisedRoute.POST("/routing-configurations/:provider", configuration.SaveConfigurationHandler)
	authorisedRoute.Any("/receive/:provider", receive.USSDReceiveHandler)

	authorisedRoute.Any("/send/:provider", send.USSDSendHandler)

	// Server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
