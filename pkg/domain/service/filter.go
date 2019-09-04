package domain

import "github.com/mpppk/sutaba-server/pkg/domain/model"

func IsTargetMessage(user *model.User, message *model.Message, keyword string) string {
	_, ok := message.GetFirstMedia()
	if !ok {
		return "message is ignored because it has no media"
	}

	if !message.HasKeyWord(keyword) {
		return "message is ignored because it has no keyword"
	}

	if user.IsOwnMessage(message) {
		return "message is ignored because it is sent by bot"
	}
	return ""
}
