package entities

const STATUS_SUCCESS = "success"
const STATUS_PENDING = "pending"
const STATUS_QUEUED = "queued"
const STATUS_FAILED = "failed"

type GenericUSSDReceivePayload struct {
	Provider      string `json:"provider"`
	MessageType   string `json:"messageType"`
	SpId          string `json:"spId"`
	ServiceId     string `json:"serviceId"`
	SessionId     string `json:"sessionId"`
	Network       string `json:"network"`
	Msisdn        string `json:"msisdn"`
	Reference     string `json:"reference"`
	AccessCode    string `json:"accessCode"`
	ServiceCode    string `json:"serviceCode"`
	AccessString  string `json:"accessString"`
	USSDString    string `json:"ussdString"`
	Timestamp     string `json:"timestamp"`
	OperationType string `json:"operationType"`
	CodeScheme    string `json:"codeScheme"`
}



type SendUSSDGenericRequest struct {
	AccessCode    string `json:"accessCode"`
	Network       string `json:"network"`
	SessionId     string `json:"sessionId"`
}