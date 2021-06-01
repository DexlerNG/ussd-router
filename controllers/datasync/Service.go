package datasync

import (
	"encoding/json"
	"fmt"
	"subscription-billing-engine/controllers/datasync/implementation"
	"subscription-billing-engine/entities"
	"subscription-billing-engine/startups/queues"
)

var providersMapping = map[string]NetworkDataSyncInterface{
	"exchange": &implementation.ExchangeProviderDataSync{},
}

func GetProvider(provider string) NetworkDataSyncInterface {
	return providersMapping[provider]
}


func ProcessDataSync(payload []byte) error {

	queuePayload := map[string]string{}
	err := json.Unmarshal(payload, &queuePayload)
	fmt.Println("Data Stream", queuePayload, "error", err,"PP",queuePayload["provider"])
	if err != nil {
		return err
	}

	providersInterface := GetProvider(queuePayload["provider"])
	if providersInterface == nil{
		fmt.Println("Provider Not Supported")
		return nil
	}

	err, response := providersInterface.ProcessDataSync(queuePayload["data"])
	if err != nil {
		return err
	}

	response.Provider = queuePayload["provider"];
	var eventName string

	if response.Mode == "subscription"{
		eventName = "datasync.subscription"
	}else{
		eventName = "datasync.unsubscription"

	}

	err = queues.RabbitMQPublishToExchange(queues.DEFAULT_EXCHANGE, entities.EventPayload{
		Event: eventName,
		Data: response,
	})
	if err != nil {
		return err
	}
	fmt.Println("Response", response)
	return nil
}


