package subscribe

import (
	"encoding/json"
	"fmt"
	"subscription-billing-engine/controllers/subscribe/implementation"
)

var providersMapping = map[string]NetworkSubscribeInterface{
	"exchange": &implementation.ExchangeProviderRequest{},
}

func GetProvider(provider string) NetworkSubscribeInterface {
	return providersMapping[provider]
}

func ProcessInitiateSubscription(payload []byte) error {

	queuePayload := map[string]interface{}{}
	err := json.Unmarshal(payload, &queuePayload)
	fmt.Println("Data Stream", queuePayload, "error", err, "PP", queuePayload["provider"].(string))
	if err != nil {
		return err
	}

	providersInterface := GetProvider(queuePayload["provider"].(string))

	//convert map string of interface to Class parameter of the interface
	if err := providersInterface.Ingest(queuePayload["data"].(map[string]interface{})); err != nil {
		return err
	}

	fmt.Println("Request", providersInterface)

	if queuePayload["mode"].(string) == "unsubscribe" {
		err = providersInterface.ProcessUnsubscription()
	} else {
		err = providersInterface.InitiateSubscription()
	}

	if err != nil {
		fmt.Println("InitiateSubscription", err)
	}
	return nil
}
