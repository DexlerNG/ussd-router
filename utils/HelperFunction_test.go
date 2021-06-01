package utils

import (
	"fmt"
	"testing"
)

func TestConvertMapToMongoDBQuery(t *testing.T) {
	queryMap := map[string]string{
		"page": "1",
		"limit": "50",
		"fieldFilter1": "client",
		"sort": "name,asc",
	}
	result := ConvertMapToMongoDBQuery(queryMap)
	fmt.Println("Query", result)
}

