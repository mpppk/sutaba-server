package model

import (
	"strconv"
)

type UserID int64
type UserName string

type User struct {
	ID   UserID
	Name UserName
}

func NewTwitterUser(id int64, name string) User {
	return User{
		ID:   UserID(id),
		Name: UserName(name),
	}
}

func (u *User) GetIDStr() string {
	return strconv.FormatInt(int64(u.ID), 10)
}

func (u *User) IsOwnMessage(message *Message) bool {
	if message == nil {
		return false
	}
	return message.User.ID == u.ID
}
