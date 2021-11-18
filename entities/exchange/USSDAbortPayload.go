package exchange

import (
	"encoding/xml"
)

type USSDAbortBody struct {
	SenderCB    string ` xml:"loc:senderCB"`
	ReceiveCB   string ` xml:"loc:receiveCB"`
	AbortReason string `xml:"loc:abortReason"`
}

type USSDAbortRequestBody struct {
	AbortUSSDBody USSDAbortBody `xml:"loc:sendUssdAbort"`
}

type USSDAbortRequestHeader struct {
	XMLName xml.Name `xml:"RequestSOAPHeader"`
	SpId       string `xml:"spId"`
	SpPassword string `xml:"spPassword"`
	ServiceId  string `xml:"serviceId"`
	TimeStamp  string `xml:"timeStamp"`
}

type USSDAbortHeader struct {
	RequestSOAPHeader USSDAbortRequestHeader `json:"RequestSOAPHeader" xml:"RequestSOAPHeader"`
}

type USSDAbortPayload struct {
	XMLName xml.Name             `xml:"soapenv:Envelope"`
	XmlNS   string               `xml:"xmlns:soapenv,attr"`
	XmlNLoc string               `xml:"xmlns:loc,attr"`
	Header  USSDAbortHeader      `xml:"soapenv:Header"`
	Body    USSDAbortRequestBody `xml:"soapenv:Body"`
}
