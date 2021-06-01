package receive

import (
	"github.com/labstack/echo/v4"
	"ussd-router/entities"
)

type USSDReceiveInterface interface {
	//GetClass() struct{}
	Process([]byte) (error, *entities.GenericUSSDReceivePayload)
	ResolveClientResponse(echo.Context, []byte)  error
}
