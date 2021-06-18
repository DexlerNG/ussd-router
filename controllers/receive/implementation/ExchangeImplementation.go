package implementation

import (
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
	"ussd-router/entities"
	"ussd-router/entities/exchange"
)

var messageTypeMap = map[string]string{
	"0": "begin",
	"1": "continue",
	"2": "end",
}
var operationTypeMap = map[string]string{
	"1": "request",
	"2": "notify",
	"3": "response",
	"4": "release",
}

type ExchangeReceiveImplementation struct {
	Payload exchange.USSDReceivePayload `json:"payload"`
}

func (request *ExchangeReceiveImplementation) Validate() error {
	return nil
	//return validation.ValidateStruct(request,
	//	validation.Field(&request.SpId, validation.Required),
	//	validation.Field(&request.CallbackURL, is.URL, validation.Required),
	//	validation.Field(&request.SpPassword, validation.Required),
	//	validation.Field(&request.ProductName, validation.Required),
	//	validation.Field(&request.AccessCode, validation.Required),
	//	validation.Field(&request.ServiceId, validation.Required),
	//)
}

func (request *ExchangeReceiveImplementation) Process(byteData []byte) (error, *entities.GenericUSSDReceivePayload) {

	//
	ussdReceive := exchange.USSDReceivePayload{}
	if err := xml.Unmarshal(byteData, &ussdReceive); err != nil {
		// do error check
		fmt.Println(err)
		return err, nil
	}

	response := &entities.GenericUSSDReceivePayload{
		Provider:      "exchange",
		MessageType:   messageTypeMap[ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.MsgType],
		SpId:          ussdReceive.Header.NotifySoapHeader.SpId,
		ServiceId:     ussdReceive.Header.NotifySoapHeader.ServiceId,
		Timestamp:     ussdReceive.Header.NotifySoapHeader.TimeStamp,
		USSDString:    ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.USSDString,
		Msisdn:        ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.Msisdn,
		CodeScheme:    ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.CodeScheme,
		SessionId:     ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.SenderCB,
		Reference:     ussdReceive.Header.NotifySoapHeader.TraceUniqueId,
		OperationType: operationTypeMap[ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.UssdOpType],
	}
	if response.MessageType == "begin" {
		fmt.Println("Begin USSD")
		response.ServiceCode = strings.TrimPrefix(ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.USSDString, "*")
		response.ServiceCode = strings.TrimSuffix(response.ServiceCode, "#")
		response.AccessCode = response.ServiceCode
		fmt.Println("Service Code on Begin", response.ServiceCode)
	} else {
		response.AccessCode = ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.ServiceCode
		response.ServiceCode = ussdReceive.Body.USSDReceiveNotifyUSSDReceptionBody.ServiceCode
	}
	fmt.Println("USSD Receive From Exchange", ussdReceive)
	return nil, response
}

func (request *ExchangeReceiveImplementation) ResolveClientResponse(c echo.Context, byteData []byte) error {
	return c.XML(200, &exchange.USSDReceiveResponsePayload{
		XmlNS:   "http://schemas.xmlsoap.org/soap/envelope/",
		XmlNLoc: "http://www.csapi.org/schema/parlayx/ussd/notification/v1_0/local",
		Header:  exchange.USSDReceiveResponseHeader{},
		Body: exchange.USSDReceiveResponseBody{
			USSDReceiveNotifyUSSDReceptionBody: exchange.USSDReceiveNotifyUSSDdReceptionResponse{
				Result: "0",
			},
		},
	})
}
