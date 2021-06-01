package receive

import (
	"ussd-router/controllers/configuration/implementation"
)

var providersMapping = map[string]USSDReceiveInterface{
	"exchange": &implementation.ExchangeConfiguration{},
}

func GetUSSDReceiveProvider(provider string) USSDReceiveInterface {
	return providersMapping[provider]
}