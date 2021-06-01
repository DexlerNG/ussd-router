package charge

import (
	"context"
	"encoding/json"
	"fmt"
	"subscription-billing-engine/controllers/charge/implementation"
	"subscription-billing-engine/entities"
	redis "subscription-billing-engine/startups/cache"
	"subscription-billing-engine/startups/queues"
	"time"
)

var providersMapping = map[string]NetworkChargeInterface{
	"exchange": &implementation.ExchangeProviderRequest{},
}

func GetProvider(provider string) NetworkChargeInterface {
	return providersMapping[provider]
}

func ProcessGetAuthorization(payload []byte) error {

	queuePayload := map[string]interface{}{}
	err := json.Unmarshal(payload, &queuePayload)
	fmt.Println("Data Stream", queuePayload, "error", err)
	if err != nil {
		return err
	}

	providersInterface := GetProvider(queuePayload["provider"].(string))

	//convert map string of interface to Class parameter of the interface
	if err := providersInterface.Ingest(queuePayload["data"].(map[string]interface{})); err != nil {
		return err
	}

	fmt.Println("Request", providersInterface)

	err, response := providersInterface.GetAuthorization()
	fmt.Println("GetAuthorization", err, response)

	if err != nil || response.Status == entities.STATUS_FAILED {
		//update payment service of failed payment
		fmt.Println("Failed", err, &response)
		return nil
	}

	//handle Queued
	if response.Status == entities.STATUS_QUEUED {
		queuePayload["authorizationResponse"] = response
		byteData, err := json.Marshal(queuePayload["data"])
		if err != nil {
			//handle error
			//ProcessError(err, sms)
			return err

		}
		redisResponse := redis.GetRedisClient().Set(context.Background(), "sub-engine:pending-payment:"+response.Reference, byteData, time.Hour)
		fmt.Println("redis", redisResponse.String())
		if redisResponse.Err() != nil {
			return redisResponse.Err()
		}
		return nil
	}
	//handle Success

	if response.Status == entities.STATUS_SUCCESS {
		return nil
	}
	return nil
}

func ProcessCharge(payload []byte) error {

	queuePayload := map[string]string{}
	err := json.Unmarshal(payload, &queuePayload)
	fmt.Println("Data Stream", queuePayload, "error", err)
	if err != nil {
		return err
	}

	provider := queuePayload["provider"]
	providersInterface := GetProvider(provider)
	err, consentResponse := providersInterface.ParseAuthorizationConsentResponse(queuePayload["data"])
	fmt.Println("Parse Response", consentResponse, "error", err)
	if err != nil {
		return err
	}

	if !consentResponse.IsConsentGiven {
		err = queues.RabbitMQPublishToQueue("payment.events", entities.EventPayload{
			Event: "airtime.billing.callback",
			Data: map[string]interface{}{
				"provider":  queuePayload["provider"],
				"reference": queuePayload["reference"],
				"status":    entities.STATUS_FAILED,
				"error":     "Consent Not Given",
				"reason":    "Consent Not Given",
			},
		})

		if err != nil {
			return err
		}
	}
	redisResponse, err := redis.GetRedisClient().Get(context.Background(), "sub-engine:pending-payment:"+queuePayload["reference"]).Result()
	fmt.Println("redis ereResponse", redisResponse, "error", err)

	if err != nil {
		return err
	}
	err = providersInterface.IngestByte([]byte(redisResponse))
	if err != nil {
		return err
	}

	err, response := providersInterface.Charge(consentResponse)
	fmt.Println("In", response)

	if err != nil || response.Status == entities.STATUS_FAILED {
		//update payment service of failed payment
		fmt.Println("Failed", err, &response)
		err = queues.RabbitMQPublishToQueue("payment.events", entities.EventPayload{
			Event: "airtime.billing.callback",
			Data: map[string]interface{}{
				"provider":  queuePayload["provider"],
				"reference": queuePayload["reference"],
				"status":    entities.STATUS_FAILED,
				"msisdn":    response.Msisdn,
				"error":     "Insufficient Balance",
				"reason":    response.Error,
			},
		})

		if err != nil {
			return err
		}
		return nil
	}

	//handle Queued
	if response.Status == entities.STATUS_QUEUED {
		byteData, err := json.Marshal(map[string]interface{}{
			"type":            "charge",
			"reference":       queuePayload["reference"],
			"provider":        queuePayload["provider"],
			"consentResponse": consentResponse,
			"requestPayload":  redisResponse,
		})
		if err != nil {
			//handle error
			//ProcessError(err, sms)
			return err

		}
		redisSetResponse := redis.GetRedisClient().Set(context.Background(), "sub-engine:data-sync:"+response.Reference, byteData, time.Hour)
		fmt.Println("redis", redisSetResponse.String())
		if redisSetResponse.Err() != nil {
			return redisSetResponse.Err()
		}
		return nil
	}
	//handle Success

	if response.Status == entities.STATUS_SUCCESS {

		//TODO Spool event

		err = queues.RabbitMQPublishToQueue("payment.events", entities.EventPayload{
			Event: "airtime.billing.callback",
			Data: map[string]interface{}{
				"provider":  queuePayload["provider"],
				"reference": queuePayload["reference"],
				"status":    entities.STATUS_SUCCESS,
				"msisdn":    response.Msisdn,
				"amount": response.Amount,
			},
		})

		return nil
	}
	return nil
}
