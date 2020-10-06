package domain

import "github.com/mpppk/sutaba-server/pkg/domain/model"

// IsTargetMessage checks if the given message is a target.
func IsTargetMessage(user *model.User, message *model.Message) (bool, bool, string) {
	if user.IsOwnMessage(message) {
		return false, false, "message is ignored because it is sent by bot"
	}

	if message.ReferencedMessage != nil && user.IsOwnMessage(message.ReferencedMessage) {
		return false, false, "message is ignored because it refer bot message"
	}

	// メッセージが画像を含む場合は、referされたメッセージは無視して判定
	if message.MediaNum != 0 {
		return true, false, ""
	}

	if message.ReferencedMessage == nil {
		return false, false, "message is ignored because it has no media and does not have referenced message"
	}

	if message.ReferencedMessage.MediaNum != 0 {
		return false, true, ""
	}

	return false, false, "message is ignored because referred message has no media"
}
