package implementation

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
	"ussd-router/entities/exchange"
	"ussd-router/models"
	"ussd-router/services"
)

var messageTypeMap = map[string]string{
	"begin":    "0",
	"continue": "1",
	"end":      "2",
}
var operationTypeMap = map[string]string{
	"request":  "1",
	"notify":   "2",
	"response": "3",
	"release":  "4",
}

type ExchangeSendImplementation struct {
	Provider      string `json:"provider"`
	SpId          string `json:"spId"`
	SpPassword    string `json:"spPassword"`
	ServiceId     string `json:"serviceId"`
	AccessCode    string `json:"accessCode"`
	Msisdn        string `json:"msisdn"`
	SessionId     string `json:"sessionId"`
	Network string `json:"network"`
	CallbackURL   string `json:"callbackURL"`
	MessageType   string `json:"messageType"`
	USSDString    string `json:"ussdString"`
	OperationType string `json:"operationType"`
	CodeScheme    string `json:"codeScheme"`
}

func (request *ExchangeSendImplementation) IngestByte(byteData []byte) error {
	if err := json.Unmarshal(byteData, &request); err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	return nil
}
func (request *ExchangeSendImplementation) Validate() error {
	request.MessageType = strings.ToLower(request.MessageType)
	request.Network = strings.ToLower(request.Network)
	if validation.IsEmpty(request.OperationType) {
		request.OperationType = "request"
	} else {
		request.OperationType = strings.ToLower(request.OperationType)
	}

	return validation.ValidateStruct(request,
		//validation.Field(&request.SpId, validation.Required),
		//validation.Field(&request.CallbackURL, is.URL),
		//validation.Field(&request.SpPassword, validation.Required),
		validation.Field(&request.MessageType, validation.In("begin", "continue", "end"), validation.Required),
		validation.Field(&request.USSDString, validation.Required),
		validation.Field(&request.OperationType, validation.Required),
		validation.Field(&request.Network, validation.In("mtn", "airtel", "glo", "9mobile"), validation.Required),
		validation.Field(&request.CodeScheme, validation.Required),
		validation.Field(&request.AccessCode, validation.Required),
		//validation.Field(&request.ServiceId, validation.Required),
		validation.Field(&request.Msisdn, validation.Required),
	)
}

func (request *ExchangeSendImplementation) Send(config *models.RoutingConfiguration) error {

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(config.SpId + config.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)

	var soapBody []byte
	var err error
	//if request.MessageType == "end"{
	//	fmt.Println("Itt is an end request")
	//	soapRequest := exchange.USSDAbortPayload{
	//		XmlNS:   "http://schemas.xmlsoap.org/soap/envelope/",
	//		XmlNLoc: "http://www.csapi.org/schema/parlayx/ussd/send/v1_0/local",
	//		Header: exchange.USSDAbortHeader{
	//			RequestSOAPHeader: exchange.USSDAbortRequestHeader{
	//				SpId:       request.SpId,
	//				SpPassword: md5Value,
	//				ServiceId:  request.ServiceId,
	//				TimeStamp:  timestamp,
	//			},
	//		},
	//		Body: exchange.USSDAbortRequestBody{
	//			AbortUSSDBody: exchange.USSDAbortBody{
	//				SenderCB: request.SessionId,
	//				ReceiveCB: request.SessionId,
	//				AbortReason: request.USSDString,
	//			},
	//		},
	//	}
	//	soapBody, err = xml.Marshal(soapRequest)
	//	if err != nil {
	//		// do error check
	//		fmt.Println(err)
	//		return err
	//	}
	//}else {
	//
	//}

	soapRequest := exchange.USSDSendPayload{
		XmlNS:   "http://schemas.xmlsoap.org/soap/envelope/",
		XmlNLoc: "http://www.csapi.org/schema/parlayx/ussd/send/v1_0/local",
		Header: exchange.USSDSendHeader{
			RequestSOAPHeader: exchange.USSDSendRequestHeader{
				SpId:       config.SpId,
				SpPassword: md5Value,
				ServiceId:  config.ServiceId,
				TimeStamp:  timestamp,
				OA:         request.Msisdn,
				FA:         request.Msisdn,
			},
		},
		Body: exchange.USSDSendBody{
			SendUSSDBody: exchange.USSDSendUSSDBody{
				Msisdn:      request.Msisdn,
				MsgType:     messageTypeMap[request.MessageType],
				UssdOpType:  operationTypeMap[request.OperationType],
				ServiceCode: request.AccessCode,
				CodeScheme:  request.CodeScheme,
				USSDString:  request.USSDString,
			},
		},
	}

	if validation.IsEmpty(request.SessionId) && request.MessageType == "begin" {
		soapRequest.Body.SendUSSDBody.SenderCB = uuid.New().String()
		soapRequest.Body.SendUSSDBody.ReceiveCB = "0xFFFFFFFF"
	} else {
		soapRequest.Body.SendUSSDBody.SenderCB = request.SessionId
		soapRequest.Body.SendUSSDBody.ReceiveCB = request.SessionId
	}

	soapBody, err = xml.Marshal(soapRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err
	}

	//
	//soapBody, err := xml.Marshal(soapRequest)
	//if err != nil {
	//	// do error check
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))

	err, _ = services.ExchangeAggregatorSendUSSD(os.Getenv("EXCHANGE_SEND_USSD_BASE_URL"), soapBody)
	fmt.Println("response err", err)
	if err != nil {
		return err
	}
	return nil
}
