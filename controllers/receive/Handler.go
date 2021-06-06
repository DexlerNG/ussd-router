package receive

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
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

	fmt.Println("Calling Goroutine MakeHTTPCallToURL")
	go MakeHTTPCallToURL(os.Getenv("TEMPORARY_ROUTING_URL"), genericPayload)
	return providersInterface.ResolveClientResponse(c, nil)
}