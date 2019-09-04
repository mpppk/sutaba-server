package model

import (
	"strconv"
)

type User struct {
	ID   int64 // FIXME do not expose
	Name string
}

func NewTwitterUser(id int64, name string) User {
	return User{
		ID:   id,
		Name: name,
	}
}

func (u *User) GetIDStr() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *User) IsOwnMessage(message *Message) bool {
	if message == nil {
		return false
	}
	return message.User.ID == u.ID
}
