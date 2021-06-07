package configuration

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"ussd-router/models"
	"ussd-router/utils"
)


func SaveConfigurationHandler(c echo.Context) error {

	providerImplementation := GetConfigurationProvider(c.Param("provider"))
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

	if err := providerImplementation.Validate(); err != nil{
		return utils.ValidationResponse(c, err.Error())
	}

	err, config := providerImplementation.Process()
	if err != nil{
		return utils.ErrorResponse(c, err.Error())
	}
	//process Configuration

	config["provider"] = c.Param("provider")
	fmt.Println("Config", config)
	err, insertion := models.SaveConfigurationFromMap(config)
	if err != nil{
		return utils.ErrorResponse(c, err.Error())
	}

	fmt.Println("Inserted", insertion)
	return utils.AcceptedResponse(c, "Configuration Saved")
}
//func RemoveConfigurationHandler(c echo.Context) error {
//	smsRequest := requests.ConfigurationRequest{
//		ClientId: c.Get("clientId").(string),
//	}
//	if err := c.Bind(&smsRequest); err != nil {
//		return utils.ValidationResponse(c, err.Error())
//	}
//	err := smsRequest.ValidateCreateConfiguration()
//	if err != nil{
//		return utils.ValidationResponse(c, err.Error())
//	}
//
//	smsRequest.Network = strings.ToLower(smsRequest.Network)
//	smsRequest.Keyword = strings.ToLower(smsRequest.Keyword)
//
//	//check if there is a config with the same network and shortcode and keyword
//	err, result := models.FindConfigurationByShortCodeAndNetworkAndKeyword(smsRequest.ShortCode, smsRequest.Network, smsRequest.Keyword)
//	fmt.Print("Result: ", result)
//	fmt.Print("SMS Request", smsRequest, "err", err)
//	if err == nil && result.ClientId != smsRequest.ClientId {
//		return utils.ValidationResponse(c, "Shortcode and keyword on the network chosen is not available to use for this client.")
//	}
//	//save configuration
//	err, _ = models.SaveConfiguration(smsRequest)
//	return utils.AcceptedResponse(c, "Configuration Saved")
//}