package subscribe

type NetworkSubscribeInterface interface {
	//GetClass() struct{}
	Ingest(map[string]interface{}) error
	IngestByte([]byte) error
	GetReference() string
	Validate() error
	InitiateSubscription() error
	ProcessUnsubscription() error
}
