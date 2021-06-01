package entities

const STATUS_SUCCESS = "success"
const STATUS_PENDING = "pending"
const STATUS_QUEUED = "queued"
const STATUS_FAILED = "failed"

type GenericChargeResponse struct {
	Status    string      `json:"status"`
	Reference string      `json:"reference"`
	Error     string      `json:"error"`
	Network   string      `json:"network"`
	Msisdn    string      `json:"msisdn"`
	Amount    uint        `json:"amoutt"`
	Data      interface{} `json:"meta"`
}

type GenericUSSDReceivePayload struct {
	Provider      string `json:"provider"`
	MessageType   string `json:"messageType"`
	SpId          string `json:"spId"`
	ServiceId     string `json:"serviceId"`
	SessionId     string `json:"sessionId"`
	Msisdn        string `json:"msisdn"`
	Reference     string `json:"reference"`
	AccessCode    string `json:"accessCode"`
	USSDString    string `json:"ussdString"`
	Timestamp     string `json:"timestamp"`
	OperationType string `json:"operationType"`
	CodeScheme    string `json:"codeScheme"`
}
