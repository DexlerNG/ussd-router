package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ConfigurationRequest struct {
	ClientId    string `json:"clientId"`
	ShortCode   string `json:"shortcode"`
	Keyword     string `json:"keyword"`
	Network     string `json:"network"`
	CallbackURL string `json:"callbackURL"`
}

func (request ConfigurationRequest) ValidateCreateConfiguration() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.ShortCode, validation.Required),
		validation.Field(&request.CallbackURL, is.URL),
		//validation.Field(&request.Keyword, validation.Required),
		//validation.Field(&request.Network, validation.Required),
	)
}