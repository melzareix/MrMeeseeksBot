package Models

import (
	"golang.org/x/oauth2"
	"github.com/satori/go.uuid"
)


// Chat bot User Model
type User struct {
	Uuid string `json:"uuid"`
	Token *oauth2.Token `json:"token"`
}

func NewUser() *User {
	uniqueId := uuid.NewV4().String()
	return &User{Uuid: uniqueId}
}
