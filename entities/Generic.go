package entities

const STATUS_SUCCESS = "success"
const STATUS_PENDING = "pending"
const STATUS_QUEUED = "queued"
const STATUS_FAILED = "failed"

type ChargeAuthorizationResponse struct {
	//status can be success(call charge immediately), pending (), queued (It has been queued and we are waiting for feedback), failed( just failed)
	Status    string      `json:"status"`
	Reference string      `json:"reference"`
	Error     string      `json:"error"`
	Data      interface{} `json:"data"`
}

type GenericConsentResponse struct {
	Network        string      `json:"network"`
	Reference      string      `json:"reference"`
	IsConsentGiven bool        `json:"isConsentGiven"`
	Msisdn         string      `json:"msisdn"`
	AuthToken      string      `json:"authToken"`
	Meta           interface{} `json:"meta"`
}

type GenericChargeResponse struct {
	Status    string      `json:"status"`
	Reference string      `json:"reference"`
	Error     string      `json:"error"`
	Network   string      `json:"network"`
	Msisdn    string      `json:"msisdn"`
	Amount    uint        `json:"amoutt"`
	Data      interface{} `json:"meta"`
}

type GenericDataSyncPayload struct {
	Provider  string `json:"provider"`
	SpId      string `json:"spId"`
	ServiceId string `json:"serviceId"`
	ProductId string `json:"productId"`
	Msisdn    string `json:"msisdn"`
	Reference string `json:"reference"`
	Validity  uint   `json:"validity"`
	AutoRenew bool   `json:"autoRenew"`
	Mode      string `json:"mode"`
	StartTime uint   `json:"startTime"`
	EndTime   uint   `json:"endTime"`
	Meta      []Item `json:"meta"`
}
