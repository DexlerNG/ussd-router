package configuration

type USSDConfigurationInterface interface {
	//GetClass() struct{}
	IngestByte([]byte) error
	Validate() error
	Process() error
}
