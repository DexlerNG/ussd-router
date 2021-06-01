package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var soapReq = &http.Request{
	Header: map[string][]string{
		"Content-Type": {"text/xml; charset=utf-8"},
	},
}

func ExchangeAggregatorRegisterUSSDEndpoint(URL string, body []byte) (error, *[]byte) {
	//add //authorizationService/services/authorization
	//reqURL, _ := url.Parse(URL)
	reqURL, _ := url.Parse(URL + "/USSDNotificationManagerService/services/USSDNotificationManager")
	soapReq.URL = reqURL
	soapReq.Method = "POST"
	soapReq.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp, err := http.DefaultClient.Do(soapReq)
	if err != nil {
		// handle error
		fmt.Println("Unmarshal Error: ", err.Error())
		return err, nil
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//call slack and drop Message
		fmt.Println("Unmarshal Error: ", err.Error())
		return err, nil
	}
	fmt.Println("Response Body: ", string(body))
	//tto test the response
	return nil, &body
}
