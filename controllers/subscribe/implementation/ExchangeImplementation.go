package implementation

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"subscription-billing-engine/entities/exchange"
	"subscription-billing-engine/services"
	"time"
)

var channelMap = map[string]string{
	"sms":  "2",
	"ussd": "3",
}
var autoRenewMap = map[bool]uint{
	true:  1,
	false: 0,
}

type ExchangeProviderRequest struct {
	Provider         string `json:"provider"`
	Msisdn           string `json:"msisdn"`
	ProductName      string `json:"productName"`
	Reference        string `json:"reference"`
	Amount           uint   `json:"amount"`
	ProvidersBaseURL string `json:"providersBaseURL"`
	SpId             string `json:"spId"`
	SpPassword       string `json:"spPassword"`
	ServiceId        string `json:"serviceId"`
	ProductId        string `json:"productId"`
	Network          string `json:"network"`
	Channel          string `json:"channel"`
	AutoRenew        bool   `json:"autoRenew"`
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
		validation.Field(&request.Reference, validation.Required),
		validation.Field(&request.Amount, validation.Required, validation.Min(uint(1))),
		validation.Field(&request.ProvidersBaseURL, is.URL, validation.Required),
		validation.Field(&request.SpId, validation.Required),
		validation.Field(&request.SpPassword, validation.Required),
		validation.Field(&request.ServiceId, validation.Required),
		validation.Field(&request.Network, validation.Required),
		//validation.Field(&request.ProductName, validation.Required),
		validation.Field(&request.Channel, validation.In("sms", "ussd"), validation.Required),
	)
}

func (request *ExchangeProviderRequest) GetReference() string {
	return request.Reference
}

func (request *ExchangeProviderRequest) InitiateSubscription() error {
	fmt.Println("request", request)

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(request.SpId + request.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)
	//transactionId := "0110000" + os.Getenv("EXCHANGE_IP_PROXY_OCTET") + time.Now().Format("060102150405") + "9999" + "001"
	//get the get authorization payload and send data

	subscribeRequest := exchange.SubscriptionPayload{
		XmlLoc:  "http://www.csapi.org/schema/parlayx/subscribe/manage/v1_0/local",
		XmlV2:   "http://www.huawei.com.cn/schema/common/v2_1",
		XmlSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: exchange.SubscriptionHeader{
			RequestSOAPHeader: exchange.SubscriptionSoapHeader{
				SpId:          request.SpId,
				SpPassword:    md5Value,
				TimeStamp:     timestamp,
				ServiceId:     request.ServiceId,
				TransactionId: request.Reference,
			},
		},
		Body: exchange.SubscriptionBody{
			SubscribeProductRequest: exchange.SubscribeProductRequest{
				SubscribeProductRequest: exchange.SubscribeProductReq{
					SubscriptionInfo: exchange.SubscriptionProductInfo{
						ProductId:     request.ProductId,
						OperationCode: "en",
						IsAutoExtend:  autoRenewMap[request.AutoRenew],
						ChannelID:     channelMap[request.Channel],
					},
					UserId: exchange.SubscriptionUserID{
						ID:   request.Msisdn,
						Type: "0",
					},
				},
			},
		},
	}

	soapBody, err := xml.Marshal(subscribeRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))
	fmt.Println("subscribeRequest", subscribeRequest)

	err, _ = services.ExchangeAggregatorInitiateSubscription(request.ProvidersBaseURL, soapBody)
	//fmt.Println("response", response, "err",err)

	return nil
}

func (request *ExchangeProviderRequest) ProcessUnsubscription() error {
	fmt.Println("request", request)

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(request.SpId + request.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)
	//transactionId := "0110000" + os.Getenv("EXCHANGE_IP_PROXY_OCTET") + time.Now().Format("060102150405") + "9999" + "001"
	//get the get authorization payload and send data

	subscribeRequest := exchange.UnsubscriptionPayload{
		XmlLoc:  "http://www.csapi.org/schema/parlayx/subscribe/manage/v1_0/local",
		XmlSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: exchange.UnsubscriptionHeader{
			RequestSOAPHeader: exchange.UnsubscriptionSoapHeader{
				XmlTNS:     "http://www.huawei.com.cn/schema/common/v2_1",
				SpId:       request.SpId,
				SpPassword: md5Value,
				TimeStamp:  timestamp,
			},
		},
		Body: exchange.UnsubscriptionBody{
			UnsubscribeProductRequest: exchange.UnsubscribeProductRequest{
				UnsubscribeProductRequest: exchange.UnsubscribeProductReq{
					SubscriptionInfo: exchange.SubscriptionProductInfo{
						ProductId:     request.ProductId,
						OperationCode: "en",
						IsAutoExtend:  autoRenewMap[request.AutoRenew],
						ChannelID:     channelMap[request.Channel],
					},
					UserId: exchange.SubscriptionUserID{
						ID:   request.Msisdn,
						Type: "0",
					},
				},
			},
		},
	}

	soapBody, err := xml.Marshal(subscribeRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))
	fmt.Println("subscribeRequest", subscribeRequest)

	err, _ = services.ExchangeAggregatorInitiateSubscription(request.ProvidersBaseURL, soapBody)
	//fmt.Println("response", response, "err",err)

	return nil
}
