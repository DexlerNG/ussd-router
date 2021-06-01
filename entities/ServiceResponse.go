package entities

type ServicesStandardResponse struct{
	Data map[string]interface{} `json:"data"`
	Error string `json:"error"`
}