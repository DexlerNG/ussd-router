package entities

type RequestHeader struct {
	SpId       string `json:"spId" xml:"spId"`
	SpPassword string `json:"spPassword" xml:"spPassword"`
	ServiceId  string `json:"serviceId" xml:"serviceId"`
	TimeStamp  string `json:"timeStamp" xml:"timeStamp"`
}

type Header struct {
	RequestSOAPHeader RequestHeader `json:"RequestSOAPHeader" xml:"RequestSOAPHeader"`
}

type Item struct {
	Key   string `json:"key" xml:"key"`
	Value string `json:"value" xml:"value"`
}

//type ChargePayload struct {
//	XMLName xml.Name `xml:"soapenv:Envelope" `
//	Header  Header   `json:"Header" xml:"soapenv:Header"`
//	Body    Body     `json:"Body" xml:"soapenv:Body"`
//}
