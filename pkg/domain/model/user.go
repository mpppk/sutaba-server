package model

import (
	"strconv"
)

type TwitterUser struct {
	ID         int64 // FIXME do not expose
	ScreenName string
}

func NewTwitterUser(id int64, name string) TwitterUser {
	return TwitterUser{
		ID:         id,
		ScreenName: name,
	}
}

func (t *TwitterUser) GetIDStr() string {
	return strconv.FormatInt(t.ID, 10)
}
