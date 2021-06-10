package receive

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"time"
	"ussd-router/models"
	redis "ussd-router/startups/cache"
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

	if validation.IsEmpty(genericPayload.SessionId){
		fmt.Println("Empty session", genericPayload)

		return providersInterface.ResolveClientResponse(c, nil)
	}

	// check cache for URL, this makes it fast
	URL, err := redis.GetRedisClient().Get(context.Background(), "ussd-sessionId:" + genericPayload.SessionId).Result()
	fmt.Println("URL From Cache", URL, "Error", err)
	if err == nil && !validation.IsEmpty(URL){
		go MakeHTTPCallToURL(URL, genericPayload)
		return providersInterface.ResolveClientResponse(c, nil)
	}
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
	redisCacheResponse := redis.GetRedisClient().Set(context.Background(), "ussd-sessionId:" + genericPayload.SessionId, config.CallbackURL, 2 * time.Minute).String()
	fmt.Println("Set Cache SessionId - > URL", redisCacheResponse)
	go MakeHTTPCallToURL(config.CallbackURL, genericPayload)
	return providersInterface.ResolveClientResponse(c, nil)
}