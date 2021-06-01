package send

type USSDSendInterface interface {
	//GetClass() struct{}
	IngestByte([]byte) error
	Validate() error
	Send() error
}
