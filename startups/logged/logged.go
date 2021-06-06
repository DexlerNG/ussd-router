package logged

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"time"
	"ussd-router/startups/queues"
)

func New() *RequestLogger {
	return new(RequestLogger)
}

type RequestLogger struct {
	QueueName   string
	SilentError bool
	IndexName   string
}

func (logger *RequestLogger) HandleEchoLogger(next echo.HandlerFunc)  echo.HandlerFunc {
	return func(c echo.Context) error {
		var startTime = time.Now()
		if err := next(c); err != nil{
			c.Error(err)
		}


		go func() {
			reqBody := []byte{}
			var err error
			if c.Request().Body != nil { // Read
				reqBody, err = ioutil.ReadAll(c.Request().Body)
				fmt.Println("err", err)

			}

			fmt.Println("Body", reqBody)

			resBody := new(bytes.Buffer)
			_ = io.MultiWriter(c.Response().Writer, resBody)

			logData := map[string]interface{}{
				"host": c.Request().Host,
				"headers": c.Request().Header,
				"method": c.Request().Method,
				"ip": c.RealIP(),
				"url": c.Request().URL,
				"requestBody": string(reqBody),
				"responseBody": resBody.Bytes(),
				"responseStatusCode": c.Response().Status,
				"path": c.Request().URL.Path,
				"processingTime": time.Since(startTime),
			}

			jsonValue, _ := json.Marshal(logData)
			fmt.Println("Request Body", string(jsonValue))
			//err := queues.RabbitMQPublishToQueue(logger.QueueName, map[string]interface{}{
			//	"index": logger.IndexName,
			//	"data": logData,
			//})
			//if err != nil {
			//	fmt.Println("logged write error",err)
			//}
		}()

		return nil
	}
}

func (logger *RequestLogger) HandleEchoBodyDump(c echo.Context, reqBody, resBody []byte) {
	//var startTime = time.Now()

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
			//"processingTime":     time.Since(startTime),
		}

		jsonValue, _ := json.Marshal(logData)
		fmt.Println("Request Body", string(jsonValue))
		err := queues.RabbitMQPublishToQueue(logger.QueueName, logData)
		if err != nil {
			fmt.Println("logged write error",err)
		}
	}()
}
