package models


type User struct {
	Id              string            `json:"userId"`
	RoleId          string            `json:"roleId"`
	Email           string            `json:"email"`
	Profile         Profile           `json:"profile"`
	PhoneNumber     string            `json:"phoneNumber"`
	//EmailVerifiedAt int16             `json:"emailVerifiedAt"`
	//PhoneVerifiedAt int16             `json:"phoneVerifiedAt"`
	Meta            map[string]interface{} `json:"meta"`
}

type Profile struct {
	Id     string            `json:"id"`
	Name   string            `json:"name"`
	Avatar string            `json:"avatar"`
	Meta   map[string]string `json:"meta"`
}
