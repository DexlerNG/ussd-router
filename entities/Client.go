package entities

type Client struct {
	ClientId           string                 `json:"clientId" bson:"clientId"`
	Name               string                 `json:"name" bson:"name" validate:"required"`
	CountryCode        string                 `json:"countryCode" bson:"countryCode"`
	DefaultUserId      string                 `json:"defaultUserId" bson:"defaultUserId"`
	SubDomain          string                 `json:"subdomain" bson:"subdomain"`
	Logo               string                 `json:"logo" bson:"logo"`
	AppURL             string                 `json:"appURL" bson:"appURL"`
	Access             string                 `json:"access" bson:"access"`
}
