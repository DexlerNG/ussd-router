package exchange

import "encoding/xml"

type USSDReceiveNotifyUSSDdReceptionResponse struct {
	Result     string `xml:"loc:result"`
}

type USSDReceiveResponseBody struct {
	USSDReceiveNotifyUSSDReceptionBody USSDReceiveNotifyUSSDdReceptionResponse `xml:"loc:notifyUssdReceptionResponse"`
}
type USSDReceiveResponseHeader struct {
}

type USSDReceiveResponsePayload struct {
	XMLName xml.Name                  `xml:"soapenv:Envelope"`
	XmlNS   string                    `xml:"xmlns:soapenv,attr"`
	XmlNLoc string                    `xml:"xmlns:loc,attr"`
	Header  USSDReceiveResponseHeader `xml:"soapenv:Header"`
	Body    USSDReceiveResponseBody           `xml:"soapenv:Body"`
}

//`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://www.csapi.org/schema/parlayx/ussd/notification/v1_0/local">
//   <soapenv:Header/>
//   <soapenv:Body>
//			<loc:notifyUssdReceptionResponse>
//					<loc:result>0</loc:result>
//			</loc:notifyUssdReceptionResponse>
//	</soapenv:Body>
//</soapenv:Envelope>`
