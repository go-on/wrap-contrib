package model

import (
// "time"
)

var users = map[string]*User{
	"donald@duck.com": {"Donald", "entenhausen", "donald@duck.com"},
	"daisy@duck.com":  {"Daisy", "blümchen", "daisy@duck.com"},
}

type User struct {
	Name, Password, EMail string
}

func FindUser(email string) *User {
	// println("find", email)
	// time.Sleep(2 * time.Second)
	u, has := users[email]
	if !has {
		return nil
	}
	return u
}
