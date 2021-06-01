package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"ussd-router/entities/exchange"
)

var soapReq = &http.Request{
	Header: map[string][]string{
		"Content-Type": {"text/xml; charset=utf-8"},
	},
}

func ExchangeAggregatorRegisterUSSDEndpoint(URL string, body []byte) (error, *exchange.AuthorizationResponse) {
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
	//body = []byte("<soapenv:Envelope\n\txmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"\n\txmlns:v1=\"http://www.csapi.org/schema/authorization/local/v1_0\">\n\t<soapenv:Header/>\n\t<soapenv:Body>\n\t\t<v1:authorizationResponse>\n\t\t\t<v1:result>\n\t\t\t\t<resultCode>0</resultCode>\n\t\t\t\t<resultDescription>Pending</resultDescription>\n\t\t\t\t<token>123456</token>\n\t\t\t</v1:result>\n\t\t\t<v1:extensionInfo/>\n\t\t</v1:authorizationResponse>\n\t</soapenv:Body>\n</soapenv:Envelope>")

	//fmt.Println("Response Body: ", string(body))

	response := exchange.AuthorizationResponse{}
	//response := map[string]interface{}{}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		//call slack and drop Message
		fmt.Println("Unmarshal Error: ", err.Error())
		return err, nil
	}
	//fmt.Println("ExchangeAuthorizationResponse", response.Body.Authorization.Result.ResultDescription)
	return nil, &response
}
func ExchangeAggregatorProcessCharge(URL string, body []byte) (error, *exchange.ChargeResponse) {
	//add //authorizationService/services/authorization
	//reqURL, _ := url.Parse(URL)
	reqURL, _ := url.Parse(URL + "/AmountChargingService/services/AmountCharging")
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
	//return nil, true
	fmt.Println("Response Body: ", string(body))
	//tto test the response
	//body = []byte("<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n  <soapenv:Body>\n    <ns1:chargeAmountResponse xmlns:ns1=\"http://www.csapi.org/schema/parlayx/payment/amount_charging/v2_1/local\"></ns1:chargeAmountResponse>\n  </soapenv:Body>\n</soapenv:Envelope>")
	//fmt.Println("Response Body: ", string(body))

	response := exchange.ChargeResponse{}
	//response := map[string]interface{}{}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		//call slack and drop Message
		fmt.Println("Unmarshal Error: ", err.Error())
		return err, nil
	}
	//fmt.Println("ExchangeAuthorizationResponse", response.Body.Authorization.Result.ResultDescription)
	return nil, &response
}
func ExchangeAggregatorInitiateSubscription(URL string, body []byte) (error, interface{}) {
	fmt.Println("URL", URL)
	reqURL, _ := url.Parse(URL + "/SubscribeManageService/services/SubscribeManage")
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
	fmt.Println("Response Body: ", string(body), "status code", resp.StatusCode, req.Header.Get("Content-Type"), resp.Header, body)
	//tto test the response
	//body = []byte("<soapenv:Envelope\n\txmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"\n\txmlns:v1=\"http://www.csapi.org/schema/authorization/local/v1_0\">\n\t<soapenv:Header/>\n\t<soapenv:Body>\n\t\t<v1:authorizationResponse>\n\t\t\t<v1:result>\n\t\t\t\t<resultCode>0</resultCode>\n\t\t\t\t<resultDescription>Pending</resultDescription>\n\t\t\t\t<token>123456</token>\n\t\t\t</v1:result>\n\t\t\t<v1:extensionInfo/>\n\t\t</v1:authorizationResponse>\n\t</soapenv:Body>\n</soapenv:Envelope>")

	//fmt.Println("Response Body: ", string(body))
	return nil, &body

	//response := exchange.AuthorizationResponse{}
	////response := map[string]interface{}{}
	//err = xml.Unmarshal(body, &response)
	//if err != nil {
	//	//call slack and drop Message
	//	fmt.Println("Unmarshal Error: ", err.Error())
	//	return err, nil
	//}
	////fmt.Println("ExchangeAuthorizationResponse", response.Body.Authorization.Result.ResultDescription)
	//return nil, &response
}

//func FindUser(userId string) (error, *models.User) {
//	var userBaseURL = os.Getenv("USER_SERVICE_URL")
//
//	reqURL, _ := url.Parse(userBaseURL + "/v1/users/" + userId)
//	req.URL = reqURL
//	req.Method = "GET"
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		// handle error
//		println("User Error: ", err.Error())
//		return err, nil
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		//call slack and drop Message
//		println("Body Reading Slack Error: ", err.Error())
//		return err, nil
//	}
//	println("Response Body: ", string(body))
//
//	userResponse := requests.ServicesStandardResponse{}
//	err = json.Unmarshal(body, &userResponse)
//	if err != nil {
//		fmt.Println("Error While decoding client: ", err)
//		return err, nil
//	}
//	if !utils.IsStringEmpty(userResponse.Error)  {
//		return errors.New(userResponse.Error), nil
//	}
//	jsonValue, _ := json.Marshal(userResponse.Data)
//	createdUser := models.User{}
//	err = json.Unmarshal(jsonValue, &createdUser)
//	if err != nil {
//		fmt.Println("Error While decoding user: ", err)
//		return err, nil
//	}
//	return nil, &createdUser
//}
