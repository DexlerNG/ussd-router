package receive

type USSDReceiveInterface interface {
	//GetClass() struct{}
	IngestByte([]byte) error
	Process() error
}
