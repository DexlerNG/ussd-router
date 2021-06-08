package receive

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"ussd-router/models"
	"ussd-router/utils"
)

//handleUserRequest
func USSDReceiveHandler(c echo.Context) error {
	providersInterface := GetUSSDReceiveProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "USSD Receive Provider '"+c.Param("provider")+"' is not supported")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return utils.ValidationResponse(c, err.Error())
	}
	//fmt.Println("BODY", string(body))

	//Spool this to another thread
	err, genericPayload := providersInterface.Process(body)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	fmt.Println("genericPayload", genericPayload)
	var config models.RoutingConfiguration


	if !utils.IsStringEmpty(c.Param("network")) {
		genericPayload.Network = c.Param("network")
	}

	if utils.IsStringEmpty(genericPayload.Network){
		//get details
		err, config = models.FindConfigurationBySpIdAndServiceIdAndAccessCodeAndAccessString(genericPayload.SpId, genericPayload.ServiceId, genericPayload.AccessCode, genericPayload.AccessString)
		if err != nil {
			fmt.Println("Unknown Error", err, genericPayload)
			return utils.ErrorResponse(c, err.Error())
		}
		genericPayload.Network = config.Network

	}else {
		//We have network...
		err, config = models.FindConfigurationByAccessCodeAndNetworkAndAccessString(genericPayload.AccessCode, genericPayload.Network, genericPayload.AccessString)
		if err != nil {
			fmt.Println("Unknown Error", err, genericPayload)
			return utils.ErrorResponse(c, err.Error())
		}
	}
	fmt.Println("Calling Goroutine MakeHTTPCallToURL")
	go MakeHTTPCallToURL(config.CallbackURL, genericPayload)
	return providersInterface.ResolveClientResponse(c, nil)
}