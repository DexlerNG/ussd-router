package send

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"ussd-router/entities"
	"ussd-router/models"
	"ussd-router/utils"
)

//handleUserRequest
func USSDSendHandler(c echo.Context) error {


	providerImplementation := GetUSSDSendProvider(c.Param("provider"))
	if providerImplementation == nil {
		return utils.ErrorResponse(c, "Provider '" + c.Param("provider") + "' is not supported")
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	fmt.Println("Body of the req", string(body))
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



	//get generic requestt
	request := entities.SendUSSDGenericRequest{
	}
	if err = json.Unmarshal(body, &request); err != nil {
		return utils.ValidationResponse(c, err.Error())
	}

	fmt.Println("Small", request)
	//
	err, config := models.FindConfigurationByAccessCodeAndNetworkAndAccessString(request.AccessCode, request.Network, "")
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}

	if err:= providerImplementation.Send(&config); err != nil{
		return utils.ErrorResponse(c, err.Error())
	}

	return utils.AcceptedResponse(c, "USSD Send Completed")
}