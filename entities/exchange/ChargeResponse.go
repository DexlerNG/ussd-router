package exchange

type ChargeAmountResponse struct {
}

type FaultResponse struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}
type ChargeResponseBody struct {
	Fault   FaultResponse      `json:"Fault" xml:"Fault"`
	ChargeAmountResponse ChargeAmountResponse `json:"chargeAmountResponse" xml:"ns1:chargeAmountResponse"`
}

type ChargeResponse struct {
	Body    ChargeResponseBody `json:"Body" xml:"Body"`
}
