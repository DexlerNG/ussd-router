package exchange

import (
	"encoding/xml"
)

type USSDSendUSSDBody struct {
	MsgType     string `xml:"loc:msgType"`
	SenderCB    string ` xml:"loc:senderCB"`
	ReceiveCB   string ` xml:"loc:receiveCB"`
	UssdOpType  string `xml:"loc:ussdOpType"`
	Msisdn      string `xml:"loc:msIsdn"`
	ServiceCode string `xml:"loc:serviceCode"`
	CodeScheme  string `xml:"loc:codeScheme"`
	USSDString  string `xml:"loc:ussdString"`
}

type USSDSendBody struct {
	SendUSSDBody USSDSendUSSDBody `xml:"loc:sendUssd"`
}

type USSDSendRequestHeader struct {
	XMLName xml.Name `xml:"RequestSOAPHeader"`
	//XmlNS      string   `xml:"xmlns:tns,attr"`
	SpId       string `xml:"spId"`
	SpPassword string `xml:"spPassword"`
	ServiceId  string `xml:"serviceId"`
	TimeStamp  string `xml:"timeStamp"`
	OA         string `xml:"OA"`
	FA         string `xml:"FA"`
}

type USSDSendHeader struct {
	RequestSOAPHeader USSDSendRequestHeader `json:"RequestSOAPHeader" xml:"RequestSOAPHeader"`
}

type USSDSendPayload struct {
	XMLName xml.Name       `xml:"soapenv:Envelope"`
	XmlNS   string         `xml:"xmlns:soapenv,attr"`
	XmlNLoc string         `xml:"xmlns:loc,attr"`
	Header  USSDSendHeader `xml:"soapenv:Header"`
	Body    USSDSendBody   `xml:"soapenv:Body"`
}
