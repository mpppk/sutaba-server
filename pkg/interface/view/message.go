package view

import "github.com/mpppk/sutaba-server/pkg/domain/model"

type MessageView interface {
	Show(text string) error
	Reply(text string, user *model.User) error
}
