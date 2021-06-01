package configuration

import (
	"ussd-router/controllers/configuration/implementation"
)

var providersMapping = map[string]USSDConfigurationInterface{
	"exchange": &implementation.ExchangeConfiguration{},
}

func GetConfigurationProvider(provider string) USSDConfigurationInterface {
	return providersMapping[provider]
}