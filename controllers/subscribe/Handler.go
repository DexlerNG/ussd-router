package subscribe

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"subscription-billing-engine/startups/queues"
	"subscription-billing-engine/utils"
)

func QueueSubscribeHandler(c echo.Context) error {

	//verify if provider is supported
	providersInterface := GetProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "Subscribe Provider '"+c.Param("provider")+"' is not supported")
	}

	if c.Param("mode") != "subscribe" && c.Param("mode") != "unsubscribe" {
		return utils.ValidationResponse(c, "Mode can either be subscribe or unsubscribe")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return utils.ValidationResponse(c, err.Error())
	}
	//convert map string of interface to Class parameter of the interface
	if err = providersInterface.IngestByte(body); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	if err = providersInterface.Validate(); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	fmt.Println("Provider", c.Param("provider"))
	err = queues.RabbitMQPublishToQueue(queues.SUBSCRIPTION_PROCESSING_QUEUE, map[string]interface{}{
		"provider": c.Param("provider"),
		"mode":     c.Param("mode"),
		"data":     providersInterface,
	})

	if err != nil {
		return utils.ErrorResponse(c, "Unable to Process Payment")
	}

	return utils.AcceptedResponse(c, "Processing...")
}

func QueueUnSubscribeHandler(c echo.Context) error {

	//verify if provider is supported
	providersInterface := GetProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "Subscribe Provider '"+c.Param("provider")+"' is not supported")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return utils.ValidationResponse(c, err.Error())
	}
	//convert map string of interface to Class parameter of the interface
	if err = providersInterface.IngestByte(body); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	if err = providersInterface.Validate(); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	fmt.Println("Provider", c.Param("provider"))
	err = queues.RabbitMQPublishToQueue(queues.SUBSCRIPTION_PROCESSING_QUEUE, map[string]interface{}{
		"provider": c.Param("provider"),
		"data":     providersInterface,
	})

	if err != nil {
		return utils.ErrorResponse(c, "Unable to Process Payment")
	}

	return utils.AcceptedResponse(c, "Processing...")
}
