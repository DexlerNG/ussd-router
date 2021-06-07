package send

import "ussd-router/models"

type USSDSendInterface interface {
	//GetClass() struct{}
	IngestByte([]byte) error
	Validate() error
	Send(config *models.RoutingConfiguration) error
}
