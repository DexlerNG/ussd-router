package exchange

import (
	"encoding/xml"
)

type SUDBodyStartUSSDNotificationReference struct {
	Endpoint      string `json:"endpoint" xml:"endpoint"`
	InterfaceName string `json:"interfaceName" xml:"interfaceName"`
	Correlator    string `json:"correlator" xml:"correlator"`
}

type SUDBodyStartUSSDNotification struct {
	Reference                   SUDBodyStartUSSDNotificationReference `json:"reference" xml:"loc:reference"`
	USSDServiceActivationNumber string                                `json:"ussdServiceActivationNumber" xml:"loc:ussdServiceActivationNumber"`
	Criteria                    string                                `json:"criteria" xml:"loc:criteria"`
}

type SUDBody struct {
	SUDBodyStartUSSDNotification SUDBodyStartUSSDNotification `json:"startUSSDNotification" xml:"loc:startUSSDNotification"`
}

type SUDRequestHeader struct {
	XMLName    xml.Name `xml:"tns:RequestSOAPHeader"`
	XmlNS      string   `xml:"xmlns:tns,attr"`
	SpId       string   `json:"spId" xml:"tns:spId"`
	SpPassword string   `json:"spPassword" xml:"tns:spPassword"`
	ServiceId  string   `json:"serviceId" xml:"tns:serviceId"`
	TimeStamp  string   `json:"timeStamp" xml:"tns:timeStamp"`
}

type SUDHeader struct {
	RequestSOAPHeader SUDRequestHeader `json:"RequestSOAPHeader" xml:"tns:RequestSOAPHeader"`
}

type SUDPayload struct {
	XMLName xml.Name     `xml:"soapenv:Envelope"`
	XmlNS   string       `xml:"xmlns:soapenv,attr"`
	XmlNLoc string       `xml:"xmlns:loc,attr"`
	Header  SUDHeader `json:"Header" xml:"soapenv:Header"`
	Body    SUDBody   `json:"Body" xml:"soapenv:Body"`
}
