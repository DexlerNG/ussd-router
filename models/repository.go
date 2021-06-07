package models

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"ussd-router/startups"
	"ussd-router/utils"
)

type fn func(cursor *mongo.Cursor) (error, interface{})

func GetPaginated(collection string, resolveFunction fn, query map[string]string) utils.PaginatedResponse {
	result := utils.ConvertMapToMongoDBQuery(query)
	fmt.Println("Result", result)
	paginatedResponse := utils.NewPaginatedResponse(result.Page, result.Limit)
	var models []interface{}
	//Create a handle to the respective collection in the database.
	//Perform Find operation & validate against the error.
	count, countError := GetCollection(collection).CountDocuments(context.TODO(), result.Filter)
	if countError != nil || count == 0 {
		log.Print("Count Error: ", countError, "count", count)
		return paginatedResponse
	}

	log.Println("Count: ", count)

	cur, findError := GetCollection(collection).Find(context.TODO(), result.Filter, &result.Options)
	if findError != nil {
		log.Println("findError", findError)
		return paginatedResponse
	}

	//Map result to slice
	for cur.Next(context.TODO()) {
		err, model := resolveFunction(cur)
		if err != nil {
			log.Println("Cursor Error", err)
			return paginatedResponse
		}
		models = append(models, model)
	}
	// once exhausted, close the cursor
	_ = cur.Close(context.TODO())
	paginatedResponse.Limit = result.Limit
	paginatedResponse.Total = count
	paginatedResponse.Page = result.Page
	paginatedResponse.Pages = int(math.Ceil(float64(count) / float64(result.Limit)))

	inrec, _ := json.Marshal(models)
	json.Unmarshal(inrec, &paginatedResponse.Data)
	return paginatedResponse
}

func Create(collection string, data interface{})(error, string){
		insertedId, err := GetCollection(collection).InsertOne(context.TODO(), data)
		if err != nil{
			return err, ""
		}
		return nil, insertedId.InsertedID.(string)
}

func Upsert(collection string, query map[string]string, data interface{})(error, interface{}){
	//create indexes
	filter := utils.ConvertMapToMongoDBQuery(query)

	opts := options.Update().SetUpsert(true)


	//return nil, client
	res, err := GetCollection(collection).UpdateOne(context.TODO(), filter.Filter, bson.D{
		{"$set", data},
	}, opts)

	fmt.Println("Response", res, err)

	if err != nil {
		return err, nil
	}
	fmt.Println("Response", res)
	return nil, &res
}
func FindOne(collection string, query map[string]string, model interface{}) (error, interface{}){
	filter := utils.ConvertMapToMongoDBQuery(query)
	fmt.Println("Filter", filter)
	//Perform FindOne operation & validate against the error.
	err := GetCollection(collection).FindOne(context.TODO(), filter.Filter).Decode(&model)
	fmt.Println("model", model)
	if err != nil {
		return err, model
	}
	return nil, model
}
func GetCollection(collection string) *mongo.Collection {
	return startups.GetMongoDatabase().Collection(collection)
}
