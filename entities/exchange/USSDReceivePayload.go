package exchange


type USSDReceiveNotifyUSSDReceptionBody struct {
	MsgType     string `xml:"msgType"`
	SenderCB    string ` xml:"senderCB"`
	ReceiveCB   string ` xml:"receiveCB"`
	UssdOpType  string `xml:"ussdOpType"`
	Msisdn      string `xml:"msIsdn"`
	ServiceCode string `xml:"serviceCode"`
	CodeScheme  string `xml:"codeScheme"`
	USSDString  string `xml:"ussdString"`
}

type USSDReceiveBody struct {
	USSDReceiveNotifyUSSDReceptionBody USSDReceiveNotifyUSSDReceptionBody `xml:"notifyUssdReception"`
}

type USSDReceiveNotifySoapHeader struct {
	SpId          string `json:"spId" xml:"spId"`
	SpRevId       string `json:"spRevId" xml:"spRevId"`
	ServiceId     string `json:"serviceId" xml:"serviceId"`
	TimeStamp     string `json:"timeStamp" xml:"timeStamp"`
	TraceUniqueId string `json:"traceUniqueID" xml:"traceUniqueID"`
	OperatorID string `json:"OperatorID" xml:"OperatorID"` 
}

type USSDReceiveHeader struct {
	NotifySoapHeader USSDReceiveNotifySoapHeader `xml:"NotifySOAPHeader"`
}

type USSDReceivePayload struct {
	Header  USSDReceiveHeader `xml:"Header"`
	Body    USSDReceiveBody   `xml:"Body"`
}
