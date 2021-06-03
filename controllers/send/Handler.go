package send

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"ussd-router/utils"
)

//handleUserRequest
func USSDSendHandler(c echo.Context) error {
	providerImplementation := GetUSSDSendProvider(c.Param("provider"))
	if providerImplementation == nil {
		return utils.ErrorResponse(c, "Provider '" + c.Param("provider") + "' is not supported")
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	if err := providerImplementation.IngestByte(body); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	fmt.Println("Send String Request Body", string(body))
	if err := providerImplementation.Validate(); err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	if err:= providerImplementation.Send(); err != nil{
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.AcceptedResponse(c, "Content Sent")
}