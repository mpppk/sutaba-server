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

func (t *User) GetIDStr() string {
	return strconv.FormatInt(t.ID, 10)
}
