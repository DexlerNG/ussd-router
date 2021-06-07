package configuration

type USSDConfigurationInterface interface {
	//GetClass() struct{}
	IngestByte([]byte) error
	Validate() error
	Process() (error, map[string]string)
}
