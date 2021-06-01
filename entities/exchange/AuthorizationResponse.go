package exchange

import (
	"encoding/xml"
)

type AuthorizationResponse struct {
	//XMLName xml.Name     `xml:"soapenv:Envelope" `
	//XMLName xml.Name      `xml:"http://schemas.xmlsoap.org/soap/envelope/ soapenv:Envelope"`
	Body *AuthorizationResponseBody `json:"Body" xml:"Body"`
}

type AuthorizationResponseBody struct {
	Authorization AuthorizationResponseResultBody `json:"authorizationResponse" xml:"authorizationResponse"`
}

type AuthorizationResponseResultBody struct {
	Result AuthorizationResultResponse `json:"result" xml:"result"`
}
type AuthorizationResultResponse struct {
	ResultCode        int    `json:"resultCode" xml:"resultCode"`
	ResultDescription string `json:"resultDescription" xml:"resultDescription"`
	Token             string `json:"token" xml:"token"`
}

type AuthorizationCallbackBody struct {
	NotifySubscriberConsentResultResponse AuthorizationCallbackBodyNotifySubscriberConsentResultResponse `xml:"loc:notifySubscriberConsentResultResponse"`
}

type AuthorizationCallbackBodyNotifySubscriberConsentResultResponse struct {
	Result            string `json:"result" xml:"loc:result"`
	ResultDescription string `json:"resultDescription" xml:"loc:resultDescription"`
}

type AuthorizationCallback struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ soapenv:Envelope"`
	//XMLNameV2 xml.Name      `xml:"http://www.huawei.com.cn/schema/common/v2_1 xmlns:v2"`
	//XMLNameV1 xml.Name      `xml:"http://www.csapi.org/schema/authorization/local/v1_0 xmlns:v1"`
	//Header entities.Header           `json:"Header" xml:"soapenv:Header"`
	Body   AuthorizationCallbackBody `json:"Body" xml:"soapenv:Body"`
}
