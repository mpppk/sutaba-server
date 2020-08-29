package domain

import "github.com/mpppk/sutaba-server/pkg/domain/model"

func IsTargetMessage(user *model.User, message *model.Message, keyword string) string {
	if message.MediaNum == 0 {
		return "message is ignored because it has no media"
	}

	if user.IsOwnMessage(message) {
		return "message is ignored because it is sent by bot"
	}
	return ""
}
