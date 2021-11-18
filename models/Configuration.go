package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ussd-router/lib/cache"
	"ussd-router/utils"
)

var collection = "routing-configuration"

//- Incoming SMS Router
//    - shortcode
//    - network
//    - Keyword
//    - callbackURL
type RoutingConfiguration struct {
	//Id          primitive.ObjectID `json:"id" bson:"_id"`
	ClientId     string `json:"clientId" bson:"clientId"`
	SpId         string `json:"spId" bson:"spId"`
	SpPassword   string `json:"spPassword" bson:"spPassword"`
	AccessCode   string `json:"shortcode" bson:"accessCode"`
	ServiceId    string `json:"serviceId" bson:"serviceId"`
	Network      string `json:"network" bson:"network"`
	AccessString string `json:"accessString" bson:"accessString"`
	CallbackURL  string `json:"callbackURL" bson:"callbackURL"`
}

func GetCacheKey(accessCode string, network string, accessString string) string {
	return "ussd-routing-configuration:" + accessCode + ":" + network + ":" + accessString
}
func GetCacheKeyPlus(spId string, serviceId string, accessCode string, accessString string) string {
	return "ussd-routing-configuration:" + spId + ":" + serviceId + ":" + accessCode + ":" + accessString
}

func FindConfigurationByAccessCodeAndNetworkAndAccessString(accessCode string, network string, accessString string) (error, *RoutingConfiguration) {

	redisResponse, err := redis.GetRedisClient().Get(context.Background(), GetCacheKey(accessCode, network, accessString)).Result()
	fmt.Println("From Cache", redisResponse, err)

	if err == nil && redisResponse != "" {
		result := RoutingConfiguration{}
		_ = json.Unmarshal([]byte(redisResponse), &result)
		return nil, &result
	}

	fmt.Println("Could Not Get Config from cache", err)
	query := map[string]string{
		"accessCode":   accessCode,
		"network":      network,
		"accessString": accessString,
	}

	filter := utils.ConvertMapToMongoDBQuery(query)
	fmt.Print("Filter", filter.Filter)
	result := RoutingConfiguration{}
	err = GetCollection(collection).FindOne(context.TODO(), filter.Filter).Decode(&result)
	//fmt.Println("Error", err)
	if err != nil {
		return err, &result
	}

	jsonValue, _ := json.Marshal(result)
	redisResponse, err = redis.GetRedisClient().Set(context.Background(), GetCacheKey(accessCode, network, accessString), jsonValue, 0).Result()
	fmt.Println("Save Config In Redis", redisResponse, "err", err)
	return nil, &result
	//
	//err, interfaceResult := FindOne(collection, query, RoutingConfiguration{})
	//fmt.Println("interfaceResult", interfaceResult)
	//return err, i
}
func FindConfigurationBySessionIdAndNetwork(network string, sessionId string) (error, *RoutingConfiguration) {

	redisResponse, err := redis.GetRedisClient().Get(context.Background(), fmt.Sprintf("ussd-session-config:%s:%s", network, sessionId)).Result()
	fmt.Println("From Session Cache", redisResponse, err)
	if err == nil && redisResponse != "" {
		result := RoutingConfiguration{}
		_ = json.Unmarshal([]byte(redisResponse), &result)
		return nil, &result
	}

	return nil, nil

}
func FindConfigurationBySpIdAndServiceIdAndAccessCodeAndAccessString(spId string, serviceId string, accessCode string, accessString string) (error, *RoutingConfiguration) {

	redisResponse, err := redis.GetRedisClient().Get(context.Background(), GetCacheKeyPlus(spId, serviceId, accessCode, accessString)).Result()
	fmt.Println("From Cache", redisResponse, err)

	if err == nil && redisResponse != "" {
		result := RoutingConfiguration{}
		_ = json.Unmarshal([]byte(redisResponse), &result)
		return nil, &result
	}

	fmt.Println("Could Not Get Config from cache", err)
	query := map[string]string{
		"spId":         spId,
		"serviceId":    serviceId,
		"accessCode":   accessCode,
		"accessString": accessString,
	}

	filter := utils.ConvertMapToMongoDBQuery(query)
	fmt.Print("Filter", filter.Filter)
	result := RoutingConfiguration{}
	err = GetCollection(collection).FindOne(context.TODO(), filter.Filter).Decode(&result)
	//fmt.Println("Error", err)
	if err != nil {
		return err, &result
	}

	jsonValue, _ := json.Marshal(result)
	redisResponse, err = redis.GetRedisClient().Set(context.Background(), GetCacheKeyPlus(spId, serviceId, accessCode, accessString), jsonValue, 0).Result()
	fmt.Println("Save Config In Redis", redisResponse, "err", err)
	return nil, &result
}

func SaveConfiguration(request RoutingConfiguration) (error, *interface{}) {

	query := map[string]string{
		"clientId":     request.ClientId,
		"accessCode":   request.AccessCode,
		"network":      request.Network,
		"accessString": request.AccessString,
	}

	//request.CreatedAt = time.Now()
	//request.UpdatedAt = time.Now()

	err, interfaceResult := Upsert(collection, query, &request)

	redisResponse, err := redis.GetRedisClient().Del(context.Background(), GetCacheKey(request.AccessCode, request.Network, request.AccessString)).Result()
	fmt.Println("Delete Config From Redis", redisResponse, "err", err)
	return err, &interfaceResult
}

func SaveConfigurationFromMap(request map[string]string) (error, *interface{}) {

	query := map[string]string{
		//"clientId":     request.ClientId,
		"accessCode":   request["accessCode"],
		"network":      request["network"],
		"accessString": request["accessString"],
	}

	//request.CreatedAt = time.Now()
	request["updatedAt"] = time.Now().String()

	err, interfaceResult := Upsert(collection, query, &request)

	redisResponse, err := redis.GetRedisClient().Del(context.Background(), GetCacheKey(query["accessCode"], query["network"], query["accessString"])).Result()
	fmt.Println("Delete Config From Redis", redisResponse, "err", err)

	return err, &interfaceResult
}