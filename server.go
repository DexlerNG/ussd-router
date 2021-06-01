package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"ussd-router/controllers"
	"ussd-router/controllers/configuration"
	"ussd-router/controllers/receive"
	"ussd-router/utils"
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
		fmt.Println("Authenticating")
		if !utils.IsStringEmpty(c.QueryParam("clientId")) {
			c.Set("clientId", c.QueryParam("clientId"))
		} else {
			c.Set("clientId", c.Request().Header.Get("client-id"))
		}

		return next(c)
	}
}

func main() {
	err := godotenv.Load()
	fmt.Println("Error", err)
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

	//redis.GetRedisClient()
	//queues.GetRabbitMQClient()
	//defer queues.GetRabbitMQConnection().Close()

	//go queues.Listen(map[string]queues.Fn{
	//	queues.CHARGE_AUTHORIZATION_PROCESSING_QUEUE: charge.ProcessGetAuthorization,
	//	queues.CHARGE_PROCESSING_QUEUE:               charge.ProcessCharge,
	//	queues.SUBSCRIPTION_PROCESSING_QUEUE:         subscribe.ProcessInitiateSubscription,
	//	queues.DATA_SYNC_PROCESSING_QUEUE:            datasync.ProcessDataSync,
	//})
	e.Any("/health-check", controllers.HealthCheckHandler)

	//e.Use()

	authorisedRoute := e.Group("/v1")
	//authorisedRoute.Use(authenticateMiddleware)
	authorisedRoute.POST("/routing-configurations/:provider", configuration.SaveConfigurationHandler)
	//authorisedRoute.DELETE("/receive-configurations/:provider/:accessCode", configuration.SaveConfigurationHandler)
	authorisedRoute.Any("/receive/:provider", receive.USSDReceiveHandler)
	//authorisedRoute.Any("/send/:provider", receive.RouteUSSDSendHandler)


	//authorisedRoute.POST("/charge/:provider/authorization-callback", charge.AuthorizationCallbackHandler)
	//authorisedRoute.POST("/charge/:provider/authorization-callback/:reference", charge.AuthorizationCallbackHandler)
	//
	////subscription/unsubscribe to a product
	//authorisedRoute.POST("/products/:mode/:provider", subscribe.QueueSubscribeHandler)
	//
	//authorisedRoute.POST("/data-sync/:provider", datasync.QueueDataSyncHandler)

	// Server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
