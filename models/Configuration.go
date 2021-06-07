package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	redis "ussd-router/startups/cache"
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

func FindConfigurationByAccessCodeAndNetworkAndAccessString(accessCode string, network string, accessString string) (error, RoutingConfiguration) {

	redisResponse, err := redis.GetRedisClient().Get(context.Background(), GetCacheKey(accessCode, network, accessString)).Result()
	fmt.Println("From Cache", redisResponse, err)

	if err == nil && redisResponse != "" {
		result := RoutingConfiguration{}
		_ = json.Unmarshal([]byte(redisResponse), &result)
		return nil, result
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
		return err, result
	}

	jsonValue, _ := json.Marshal(result)
	redisResponse, err = redis.GetRedisClient().Set(context.Background(), GetCacheKey(accessCode, network, accessString), jsonValue, 0).Result()
	fmt.Println("Save Config In Redis", redisResponse, "err", err)
	return nil, result
	//
	//err, interfaceResult := FindOne(collection, query, RoutingConfiguration{})
	//fmt.Println("interfaceResult", interfaceResult)
	//return err, i
}
func FindConfigurationBySpIdAndServiceIdAndAccessCodeAndAccessString(spId string, serviceId string, accessCode string, accessString string) (error, RoutingConfiguration) {

	redisResponse, err := redis.GetRedisClient().Get(context.Background(), GetCacheKeyPlus(spId, serviceId, accessCode, accessString)).Result()
	fmt.Println("From Cache", redisResponse, err)

	if err == nil && redisResponse != "" {
		result := RoutingConfiguration{}
		_ = json.Unmarshal([]byte(redisResponse), &result)
		return nil, result
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
		return err, result
	}

	jsonValue, _ := json.Marshal(result)
	redisResponse, err = redis.GetRedisClient().Set(context.Background(), GetCacheKeyPlus(spId, serviceId, accessCode, accessString), jsonValue, 0).Result()
	fmt.Println("Save Config In Redis", redisResponse, "err", err)
	return nil, result
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

//
//func FindByClientId(clientId string) (error, ClientWithoutSecret) {
//	println("clientId: ", clientId)
//	//check redis
//	result := ClientWithoutSecret{}
//
//	redisResponse, _ := startups.GetRedisClient().HGet(context.TODO(), "routing", clientId).Result()
//	fmt.Println("Redis Response", redisResponse)
//	if !utils.IsStringEmpty(redisResponse) {
//		err := json.Unmarshal([]byte(redisResponse), &result)
//		if err == nil {
//			return nil, result
//		}
//	}
//
//	fmt.Println("Cache is empty, calling mongodb", redisResponse)
//
//	filter := bson.D{primitive.E{Key: "clientId", Value: clientId}}
//	//Perform FindOne operation & validate against the error.
//	err := GetCollection(collection).FindOne(context.TODO(), filter).Decode(&result)
//	if err != nil {
//		return err, result
//	}
//
//	//Return result without any error.
//	//put bank in cache
//	encodedResult, err := json.Marshal(result)
//	if err == nil {
//		hashSetResult, err := startups.GetRedisClient().HSet(context.TODO(), "routing", clientId, encodedResult).Result()
//		fmt.Println("Redis Response", hashSetResult, "Error", err)
//	}
//	result.UpdatedAt = result.UpdatedAtAsString.Unix()
//	result.CreatedAt = result.CreatedAtAsString.Unix()
//	return nil, result
//}
//
//func FindByClientIdWithSecret(clientId string) (error, Client) {
//	//check redis
//	result := Client{}
//	redisResponse, _ := startups.GetRedisClient().HGet(context.TODO(), "routing", clientId).Result()
//	fmt.Println("Redis Response", redisResponse)
//	if !utils.IsStringEmpty(redisResponse) {
//		err := json.Unmarshal([]byte(redisResponse), &result)
//		if err == nil {
//			return nil, result
//		}
//	}
//
//	fmt.Println("Cache is empty, calling mongodb", redisResponse)
//
//	filter := bson.D{primitive.E{Key: "clientId", Value: clientId}}
//	//Perform FindOne operation & validate against the error.
//	err := GetCollection(collection).FindOne(context.TODO(), filter).Decode(&result)
//	if err != nil {
//		return err, result
//	}
//
//	//Return result without any error.
//	//put bank in cache
//	encodedResult, err := json.Marshal(result)
//	if err == nil {
//		hashSetResult, err := startups.GetRedisClient().HSet(context.TODO(), "routing", clientId, encodedResult).Result()
//		fmt.Println("Redis Response", hashSetResult, "Error", err)
//	}
//	result.UpdatedAt = result.UpdatedAtAsString.Unix()
//	result.CreatedAt = result.CreatedAtAsString.Unix()
//	return nil, result
//}
//func FetchOneClient(id string) (error, ClientWithoutSecret) {
//	println("Id: ", id)
//	result := ClientWithoutSecret{}
//	//Define filter query for fetching specific document from collection
//	idObjectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return err, result
//	}
//	filter := bson.D{primitive.E{Key: "_id", Value: idObjectId}}
//	//Perform FindOne operation & validate against the error.
//	err = GetCollection(collection).FindOne(context.TODO(), filter).Decode(&result)
//	if err != nil {
//		return err, result
//	}
//	//Return result without any error.
//	return nil, result
//}
//
//func UpdateOneClient(id string, client map[string]interface{}) (error, string) {
//	//client.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
//	fmt.Println("CLient", client)
//	client["updatedAt"] = primitive.Timestamp{T: uint32(time.Now().Unix())}
//	filter := bson.M{"clientId": id} //Define filter query for fetching specific document from collection
//	updater := bson.D{primitive.E{Key: "$set", Value: client}}
//	data, err := GetCollection(collection).UpdateOne(context.TODO(), filter, updater)
//	if err != nil {
//		return err, ""
//	}
//	fmt.Println("Data", data, err)
//	hashSetResult, err := startups.GetRedisClient().HDel(context.TODO(), "routing", id).Result()
//	fmt.Println("Redis Response", hashSetResult, "Error", err)
//
//	return nil, "Client Updated Successfully"
//}
//
//func DeleteOneClient(id string) (error, string) {
//	println("Id: ", id)
//	idObjectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return err, ""
//	}
//	//client.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
//	filter := bson.D{primitive.E{Key: "_id", Value: idObjectId}} //Define filter query for fetching specific document from collection
//	data, err := GetCollection(collection).DeleteOne(context.TODO(), filter)
//	if err != nil {
//		return err, ""
//	}
//	println("Deleted Response", data)
//	return nil, "Document Deleted Successfully"
//}
