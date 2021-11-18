package utils

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

type ConversionResult struct {
	Page    int64
	Limit   int64
	Filter  bson.M
	Options options.FindOptions
}

func IsStringEmpty(val string) bool {
	return len(strings.TrimSpace(val)) == 0
}

func MapContainsString(array []string, match string) bool {
	for _, a := range array {
		if strings.ToLower(a) == strings.ToLower(match) {
			return true
		}
	}
	return false

}

func StructToMapString(class interface{}) map[string]string {
	var inInterface map[string]string
	{
	}
	mappedData, _ := json.Marshal(class)
	_ = json.Unmarshal(mappedData, &inInterface)
	return inInterface
}

func StructToMap(class interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	mappedData, _ := json.Marshal(class)
	_ = json.Unmarshal(mappedData, &inInterface)
	return inInterface
}

func ArtificialTenaryOperator(condition bool, value1 interface{}, value2 interface{}) interface{} {
	if condition {
		return value1
	}
	return value2
}

func ConvertMapToMongoDBQuery(query map[string]string) ConversionResult {

	result := ConversionResult{
		Filter:  bson.M{},
		Options: options.FindOptions{},
	}
	result.Options.SetSort(bson.M{"_id": -1})

	for k, v := range query {
		//test for page
		//query := c.QueryParams().Get("group")

		if k == "page" {
			pageConv, err := strconv.Atoi(v)
			if err != nil {
				pageConv = 1
			}
			result.Page = int64(pageConv)
			continue
		}

		if k == "limit" {
			limitConv, err := strconv.Atoi(v)
			if err != nil {
				limitConv = 1
			}
			result.Limit = int64(limitConv)
			continue
		}

		if k == "sort" {
			if IsStringEmpty(v) {
				continue
			}

			fmt.Println("Sort", v)
			//sort := bson.M{}
			splitString := strings.Split(v, ",")
			if splitString[1] == "desc" {
				result.Options.SetSort(bson.D{{splitString[0], -1}})
			} else {
				result.Options.SetSort(bson.D{{splitString[0], 1}})
			}
			continue
		}
		result.Filter[k] = v
	}

	skip := (result.Page - 1) * result.Limit
	result.Options.Skip = &skip
	result.Options.Limit = &result.Limit
	//fmt.Println("Page: ", pag e, ", Limit: ", limit)
	return result
}
