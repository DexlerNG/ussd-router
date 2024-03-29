package receive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"ussd-router/controllers/receive/implementation"
	"ussd-router/entities"
)

var providersMapping = map[string]USSDReceiveInterface{
	"exchange": &implementation.ExchangeReceiveImplementation{},
}

func GetUSSDReceiveProvider(provider string) USSDReceiveInterface {
	return providersMapping[provider]
}

func MakeHTTPCallToURL(URL string, payload *entities.GenericUSSDReceivePayload) {
	var req = &http.Request{
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=utf-8"},
		},
	}
	//Check DB or cache for URL to route, but for now, use the ENV
	reqURL, _ := url.Parse(URL)
	jsonValue, _ := json.Marshal(payload)
	fmt.Printf("Sending %s to %s", string(jsonValue), URL)
	req.URL = reqURL
	req.Method = "POST"
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonValue))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
		fmt.Println("MakeHTTPCallToURL Unmarshal Error: ", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//call slack and drop Message
		fmt.Println("MakeHTTPCallToURL Unmarshal Error: ", err.Error())
		return
	}
	//tto test the response

	fmt.Println("MakeHTTPCallToURL  Response Body: ", string(body))
	//fmt.Println("ExchangeAuthorizationResponse", response.Body.Authorization.Result.ResultDescription)
	return
}
