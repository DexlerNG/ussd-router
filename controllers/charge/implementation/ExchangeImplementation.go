package implementation

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"
	"os"
	"strconv"
	"subscription-billing-engine/entities"
	"subscription-billing-engine/entities/exchange"
	"subscription-billing-engine/services"
	"time"
)

var channelMap = map[string]string{
	"sms":  "2",
	"ussd": "3",
}

type ExchangeProviderRequest struct {
	Provider    string `json:"provider"`
	Msisdn      string `json:"msisdn"`
	ProductName string `json:"productName"`
	//ShortCode      string `json:"shortCode"`
	Reference        string `json:"reference"`
	Amount           uint   `json:"amount"`
	Description      string `json:"description"`
	ProvidersBaseURL string `json:"providersBaseURL"`
	CallbackURL      string `json:"callbackURL"`
	SpId             string `json:"spId"`
	SpPassword       string `json:"spPassword"`
	ServiceId        string `json:"serviceId"`
	Network          string `json:"network"`
	Scope            string `json:"scope"`
	Channel          string `json:"channel"`
}

func (request *ExchangeProviderRequest) IngestByte(byteData []byte) error {
	if err := json.Unmarshal(byteData, &request); err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	return nil
}

func (request *ExchangeProviderRequest) Ingest(rawData map[string]interface{}) error {
	jsonBody, err := json.Marshal(rawData)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err
	}

	if err := json.Unmarshal(jsonBody, &request); err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	return nil
}

func (request *ExchangeProviderRequest) Validate() error {

	return validation.ValidateStruct(request,
		validation.Field(&request.Msisdn, validation.Required),
		//validation.Field(&request.ShortCode, validation.Required),
		validation.Field(&request.Reference, validation.Required),
		validation.Field(&request.Amount, validation.Required, validation.Min(uint(1))),
		validation.Field(&request.ProvidersBaseURL, is.URL, validation.Required),
		validation.Field(&request.CallbackURL, is.URL),
		validation.Field(&request.SpId, validation.Required),
		validation.Field(&request.SpPassword, validation.Required),
		validation.Field(&request.ServiceId, validation.Required),
		validation.Field(&request.Network, validation.Required),
		validation.Field(&request.Scope, validation.Required),
		validation.Field(&request.ProductName, validation.Required),
		validation.Field(&request.Channel, validation.In("sms", "ussd"), validation.Required),
	)
}

func (request *ExchangeProviderRequest) GetReference() string {
	return request.Reference
}

func (request *ExchangeProviderRequest) GetAuthorization() (error, *entities.ChargeAuthorizationResponse) {

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(request.SpId + request.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)
	//transactionId := "0110000" + os.Getenv("EXCHANGE_IP_PROXY_OCTET") + time.Now().Format("060102150405") + "9999" + "001"
	//get the get authorization payload and send data
	fmt.Println("Timestamp", timestamp, "Str", md5Value)

	soapRequest := exchange.AuthorizationPayload{
		XmlNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: exchange.AuthorizationHeader{
			RequestSOAPHeader: exchange.RequestHeader{
				XmlNS:      "http://www.huawei.com.cn/schema/common/v2_1",
				SpId:       request.SpId,
				TimeStamp:  timestamp,
				SpPassword: md5Value,
			},
		},
		Body: exchange.AuthorizationBody{
			Authorization: exchange.Authorization{
				XmlNSV3:           "http://www.csapi.org/schema/authorization/local/v1_0",
				XmlNSV4:           "http://www.huawei.com.cn/schema/common/v2_1",
				EndUserIdentifier: request.Msisdn,
				//TransactionId:     transactionId,
				TransactionId: request.Reference,
				Scope:         request.Scope,
				ServiceId:     request.ServiceId,
				Amount:        request.Amount,
				Currency:      "NGN",
				Frequency:     "0",
				Description:   request.Description,
				TokenValidity: 1,
				//NotificationURL:   os.Getenv("EXCHANGE_PROXY_BASE_URL") + "/v1/charge/exchange/authorization-callback/" + request.Reference,
				//NotificationURL: "https://webhook.site/6099b649-2a94-438a-a9c5-dc1622b35f21",
				NotificationURL: os.Getenv("APP_BASE_URL") + "/v1/charge/exchange/authorization-callback/" + request.Reference,
				ExtensionInfo: []entities.Item{
					{Key: "productName", Value: request.ProductName},
					{Key: "totalAmount", Value: strconv.Itoa(int(request.Amount))},
					{Key: "channel", Value: channelMap[request.Channel]},
					{Key: "tokenType", Value: "0"},
					{Key: "serviceInterval", Value: "1"},
					{Key: "serviceIntervalUnit", Value: "1"},
				},
			},
		},
	}
	//fmt.Println("Printtln", soapRequest)
	//fmt.Printf("%+v\n", soapRequest)

	soapBody, err := xml.Marshal(soapRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err, nil
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))
	err, data := services.ExchangeAggregatorGetChargeAuthorization(request.ProvidersBaseURL, soapBody)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err, nil
	}

	response := entities.ChargeAuthorizationResponse{
		Data:      data.Body.Authorization.Result,
		Reference: request.Reference,
	}
	if data.Body.Authorization.Result.ResultCode != 0 {
		response.Status = entities.STATUS_FAILED
		response.Error = data.Body.Authorization.Result.ResultDescription
		return nil, &response
	}

	response.Status = entities.STATUS_QUEUED
	return nil, &response
}

func (request *ExchangeProviderRequest) GetAuthorizationResponse(c echo.Context) error {
	//	byteData, err := xml.Marshal(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://www.csapi.org/schema/subscriberconsnet/data/v1_0/local">
	//  <soapenv:Header/>
	//  <soapenv:Body>
	//<loc:notifySubscriberConsentResultResponse> <loc:result>0</loc:result> <loc:resultDescription>Success</loc:resultDescription>
	//</loc:notifySubscriberConsentResultResponse> </soapenv:Body>
	//</soapenv:Envelope>`)
	//	if err != nil{
	//		return err
	//	}
	return c.XML(200, exchange.AuthorizationCallback{
		Body: exchange.AuthorizationCallbackBody{
			NotifySubscriberConsentResultResponse: exchange.AuthorizationCallbackBodyNotifySubscriberConsentResultResponse{
				Result:            "0",
				ResultDescription: "Success",
			},
		},
	})
}

func (request *ExchangeProviderRequest) ParseAuthorizationConsentResponse(xmlResponse string) (error, *entities.GenericConsentResponse) {

	exchangeConsentResponse := exchange.ConsentResultResponse{}
	err := xml.Unmarshal([]byte(xmlResponse), &exchangeConsentResponse)
	fmt.Println("Error", err, "Data", exchangeConsentResponse)
	fmt.Printf("%+v\n", exchangeConsentResponse)

	if err != nil {
		return err, nil
	}
	consentResponse := &entities.GenericConsentResponse{
		IsConsentGiven: exchangeConsentResponse.Body.NotifySubscriberConsentResult.ConsentResult == 0,
		AuthToken:      exchangeConsentResponse.Body.NotifySubscriberConsentResult.AccessToken,
		Msisdn:         exchangeConsentResponse.Body.NotifySubscriberConsentResult.Subscriber.Id,
		Meta:           exchangeConsentResponse,
	}
	for _, v := range exchangeConsentResponse.Body.NotifySubscriberConsentResult.ExtensionInfo {
		if v.Key != "transactionID" {
			continue
		}

		consentResponse.Reference = v.Value
		break
	}

	if validation.IsEmpty(consentResponse.Reference) {
		consentResponse.Reference = exchangeConsentResponse.Header.NotifySoapHeader.TraceUniqueID
	}
	return nil, consentResponse
}

func (request *ExchangeProviderRequest) Charge(consent *entities.GenericConsentResponse) (error, *entities.GenericChargeResponse) {
	fmt.Printf("%+v\n", consent)
	fmt.Printf("%+v\n", request)

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(request.SpId + request.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)
	//get the get authorization payload and send data
	soapRequest := &exchange.ChargePayload{
		XmlNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: exchange.ChargeHeader{
			RequestSOAPHeader: exchange.ChargeRequestHeader{
				XmlNS:      "http://www.huawei.com.cn/schema/common/v2_1",
				SpId:       request.SpId,
				TimeStamp:  timestamp,
				SpPassword: md5Value,
				OauthToken: consent.AuthToken,
				ServiceId:  request.ServiceId,
				OA:         request.Msisdn,
				FA:         request.Msisdn,
			},
		},
		Body: exchange.ChargeBody{
			ChargeAmount: exchange.ChargeAmount{
				XmlN3:             "http://www.csapi.org/schema/parlayx/payment/amount_charging/v2_1/local",
				XmlN4:             "http://www.csapi.org/schema/parlayx/common/v2_1",
				EndUserIdentifier: request.Msisdn,
				ReferenceCode:     request.Reference + "223", //todo: replace this
				Charge: exchange.InnerCharge{
					Currency:    "NGN",
					Amount:      request.Amount,
					Description: request.Description,
					Code:        request.Reference + "225", //todo: replace this
				},
			},
		},
	}
	soapBody, err := xml.Marshal(soapRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err, nil
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))
	err, response := services.ExchangeAggregatorProcessCharge(request.ProvidersBaseURL, soapBody)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err, nil
	}
	fmt.Printf("%+v\n", response)

	chargeResponse := entities.GenericChargeResponse{
		Reference: request.Reference,
		Msisdn:    request.Msisdn,
		Network:   request.Network,
		Amount:    request.Amount,
	}

	if response.Body.Fault.FaultCode != "" {
		chargeResponse.Status = entities.STATUS_FAILED
		chargeResponse.Error = response.Body.Fault.FaultString
	} else {
		chargeResponse.Status = entities.STATUS_SUCCESS
	}
	return nil, &chargeResponse
}
