package datasync

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"subscription-billing-engine/startups/queues"
	"subscription-billing-engine/utils"
)


func QueueDataSyncHandler(c echo.Context) error {
	//verify if provider is supported
	providersInterface := GetProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "DataSync Provider '" + c.Param("provider") + "' is not supported")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	err = queues.RabbitMQPublishToQueue(queues.DATA_SYNC_PROCESSING_QUEUE, map[string]string{
		"provider": c.Param("provider"),
		"data": string(body),
	})

	if  err != nil{
		return utils.ErrorResponse(c, "Unable to Process Data Sync")
	}

	return utils.AcceptedResponse(c, "Processing...")
}