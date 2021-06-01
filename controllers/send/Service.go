package send

import (
	"ussd-router/controllers/send/implementation"
)

var providersMapping = map[string]USSDSendInterface{
	"exchange": &implementation.ExchangeSendImplementation{},
}

func GetUSSDSendProvider(provider string) USSDSendInterface {
	return providersMapping[provider]
}