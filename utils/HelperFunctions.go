package utils

import (
	"encoding/json"
	"strings"
)

func IsStringEmpty(val string) bool{
	return len(strings.TrimSpace(val)) == 0
}

func MapContainsString(array []string, match string) bool{
	for _, a := range array {
		if strings.ToLower(a) == strings.ToLower(match) {
			return true
		}
	}
	return false

}
func StructToMap(class interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	mappedData, _ := json.Marshal(class)
	json.Unmarshal(mappedData, &inInterface)
	return inInterface
}

func ArtificialTenaryOperator(condition bool, value1 interface{}, value2 interface{}) interface{}  {
	if condition{
		return value1
	}
	return value2
}