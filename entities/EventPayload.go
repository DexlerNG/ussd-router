package entities


const SMS_INBOUND_EVENT= "sms.inbound"
const ES_EVENTS_QUEUE = "elasticsearch.events"
const HTTP_WEBHOOK_EVENT = "events.http"


type EventPayload struct {
	Webhook  string      `json:"webhook"`
	Event    string      `json:"event"`
	ClientId string      `json:"clientId"`
	Data     interface{} `json:"data"`
}

type ESRunner struct {
	Index string      `json:"index"`
	Id    string      `json:"id"`
	IdKey string      `json:"idKey"`
	Data  interface{} `json:"data"`
}
