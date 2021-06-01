package requests


//- phoneNumber
//- shortCode
//- reference
//- amount
//- description
//- chargeEndpoint
//- callbackURL
//- spId
//- spPassword
//- serviceId
//- network

type ExchangeProviderRequest struct {
	Msisdn  string `json:"msisdn"`
	ShortCode string `json:"shortCode"`
	Reference string `json:"reference"`
	Amount uint `json:"amount"`
	Description string `json:"description"`
	ChargeEndpoint string `json:"chargeEndpoint"`
	CallbackURL string `json:"callbackURL"`
	SpId string `json:"spId"`
	SpPassword string `json:"spPassword"`
	ServiceId string `json:"serviceId"`
	Network string `json:"network"`
}