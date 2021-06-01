package implementation

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"os"
	"time"
	"ussd-router/entities/exchange"
	"ussd-router/services"
)

var channelMap = map[string]string{
	"sms":  "2",
	"ussd": "3",
}

type ExchangeConfiguration struct {
	SpId        string `json:"spId"`
	SpPassword  string `json:"spPassword"`
	ServiceId   string `json:"serviceId"`
	Provider    string `json:"provider"`
	ProductName string `json:"productName"`
	AccessCode  string `json:"accessCode"`
	CallbackURL string `json:"callbackURL"`
}

func (request *ExchangeConfiguration) Validate() error {
	return validation.ValidateStruct(request,
		validation.Field(&request.SpId, validation.Required),
		validation.Field(&request.CallbackURL, is.URL, validation.Required),
		validation.Field(&request.SpPassword, validation.Required),
		validation.Field(&request.ProductName, validation.Required),
		validation.Field(&request.AccessCode, validation.Required),
		validation.Field(&request.ServiceId, validation.Required),
	)
}
func (request *ExchangeConfiguration) IngestByte(byteData []byte) error {
	if err := json.Unmarshal(byteData, &request); err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	return nil
}

func (request *ExchangeConfiguration) Process() error {

	timestamp := time.Now().Format("20060102150405")
	md5Sum := md5.Sum([]byte(request.SpId + request.SpPassword + timestamp))
	md5Value := hex.EncodeToString(md5Sum[:])
	fmt.Println("Timestamp", timestamp, "Str", md5Value)

	soapRequest := exchange.SUDPayload{
		XmlNLoc: "http://www.csapi.org/schema/osg/ussd/notification_manager/v1_0/local",
		XmlNS:   "http://schemas.xmlsoap.org/soap/envelope/",
		Header: exchange.SUDHeader{
			RequestSOAPHeader: exchange.SUDRequestHeader{
				XmlNS: "http://www.huawei.com.cn/schema/common/v2_1",
				SpId:       request.SpId,
				SpPassword: md5Value,
				TimeStamp:  timestamp,
				ServiceId:  request.ServiceId,
			},
		},
		Body: exchange.SUDBody{
			SUDBodyStartUSSDNotification: exchange.SUDBodyStartUSSDNotification{
				Reference: exchange.SUDBodyStartUSSDNotificationReference{
					Endpoint: os.Getenv("APP_PROXY_BASE_URL") + "/v1/receive/exchange/receive",
					InterfaceName: request.ProductName,
					Correlator: request.AccessCode,
				},
				USSDServiceActivationNumber: request.AccessCode,
				Criteria: request.AccessCode,
			},
		},
	}

	soapBody, err := xml.Marshal(soapRequest)
	if err != nil {
		// do error check
		fmt.Println(err)
		return err
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", string(soapBody))

	err, _ = services.ExchangeAggregatorRegisterUSSDEndpoint(os.Getenv("EXCHANGE_ENDPOINT_USSD_BASE_URL"), soapBody)
	return nil
}
