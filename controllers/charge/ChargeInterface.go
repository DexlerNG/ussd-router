package charge

import (
	"github.com/labstack/echo/v4"
	"subscription-billing-engine/entities"
)

type NetworkChargeInterface interface {
	//GetClass() struct{}
	Ingest(map[string]interface{}) error
	IngestByte([]byte) error
	GetReference() string
	Validate() error
	GetAuthorization() (error, *entities.ChargeAuthorizationResponse)
	GetAuthorizationResponse(echo.Context) error
	ParseAuthorizationConsentResponse(string) (error, *entities.GenericConsentResponse)
	Charge(response *entities.GenericConsentResponse) (error, *entities.GenericChargeResponse)
}