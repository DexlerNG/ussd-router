package controllers

import (
	"github.com/labstack/echo/v4"
	"time"
	"ussd-router/utils"
)

type HealthCheckResponse struct {
	Uptime time.Duration `json:"uptime"`
	Cache  bool          `json:"cache"`
	Queue  bool          `json:"queue"`
}

var startTime = time.Now()

func uptime() time.Duration {
	return time.Since(startTime)
}
func HealthCheckHandler(c echo.Context) error {
	healthCheckResponse := HealthCheckResponse{
		//Cache: false,
		//Queue: false,
	}
	//redisResponse, err := redis.GetRedisClient().Ping(context.TODO()).Result()
	//log.Println("Redis Response", redisResponse, "Error", err)
	//if err != nil {
	//	return utils.ErrorResponse(c, err.Error())
	//}
	//healthCheckResponse.Cache = true
	//jsonPayload, _ := json.Marshal(map[string]bool{
	//	"smsIncomingRouter": true,
	//})
	//message := amqp.Publishing{
	//	Body: []byte(string(jsonPayload)),
	//}
	//err = queues.GetRabbitMQClient().Publish("health.check", "", false, false, message)
	//log.Println("QueuePayload", jsonPayload, "Error", err)
	//
	//if err != nil {
	//	return utils.ErrorResponse(c, err.Error())
	//}
	//
	//healthCheckResponse.Queue = true
	healthCheckResponse.Uptime = uptime()
	return utils.SuccessResponse(c, healthCheckResponse)

}
