package exchange

import (
	"encoding/xml"
)

type InnerCharge struct {
	Description string `json:"description" xml:"description"`
	Currency    string `json:"currency" xml:"currency"`
	Amount      uint   `json:"amount" xml:"amount"`
	Code        string `json:"code" xml:"code"`
}

type ChargeAmount struct {
	XmlN3             string      `xml:"xmlns:ns3,attr"`
	XmlN4             string      `xml:"xmlns:ns4,attr"`
	EndUserIdentifier string      `json:"endUserIdentifier" xml:"ns3:endUserIdentifier"`
	ReferenceCode     string      `json:"referenceCode" xml:"ns3:referenceCode"`
	Charge            InnerCharge `json:"charge" xml:"ns3:charge"`
}

type ChargeBody struct {
	ChargeAmount ChargeAmount `json:"chargeAmount" xml:"ns3:chargeAmount"`
}

type ChargeRequestHeader struct {
	XMLName    xml.Name `xml:"tns:RequestSOAPHeader"`
	XmlNS      string   `xml:"xmlns:tns,attr"`
	SpId       string   `json:"spId" xml:"tns:spId"`
	SpPassword string   `json:"spPassword" xml:"tns:spPassword"`
	ServiceId  string   `json:"serviceId" xml:"tns:serviceId"`
	TimeStamp  string   `json:"timeStamp" xml:"tns:timeStamp"`
	OauthToken string   `json:"oauthToken" xml:"tns:oauth_token"`
	OA         string   `xml:"tns:OA"`
	FA         string   `xml:"tns:FA"`
}

type ChargeHeader struct {
	RequestSOAPHeader ChargeRequestHeader `json:"RequestSOAPHeader" xml:"tns:RequestSOAPHeader"`
}

type ChargePayload struct {
	XMLName xml.Name     `xml:"SOAP-ENV:Envelope"`
	XmlNS   string       `xml:"xmlns:SOAP-ENV,attr"`
	Header  ChargeHeader `json:"Header" xml:"SOAP-ENV:Header"`
	Body    ChargeBody   `json:"Body" xml:"SOAP-ENV:Body"`
}