package receive

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"ussd-router/utils"
)

//handleUserRequest
func USSDReceiveHandler(c echo.Context) error{
	providersInterface := GetUSSDReceiveProvider(c.Param("provider"))
	if providersInterface == nil {
		return utils.ValidationResponse(c, "USSD Receive Provider '" + c.Param("provider") + "' is not supported")
	}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}
	fmt.Println("BODY", string(body))
	return c.XML(200, `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://www.csapi.org/schema/parlayx/ussd/notification/v1_0/local">
   <soapenv:Header/>
   <soapenv:Body>
<loc:notifyUssdReceptionResponse> <loc:result>0</loc:result>
</loc:notifyUssdReceptionResponse> </soapenv:Body>
</soapenv:Envelope>`)
}
func RouteUSSDSendHandler(c echo.Context) error{
	//providersInterface := GetProvider(c.Param("provider"))
	//if providersInterface == nil {
	//	return utils.ValidationResponse(c, "DataSync Provider '" + c.Param("provider") + "' is not supported")
	//}

	//validate request payload based on provider request data
	var body []byte
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil{
		return utils.ValidationResponse(c, err.Error())
	}

	return utils.AcceptedResponse(c, map[string]string{
		"provider": c.Param("provider"),
		"data": string(body),
	})
}



//handleClientResponse