package domain

import (
	"fmt"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

// IsTargetMessageEvent checks if the given message event is a target.
func IsTargetMessageEvent(botUser *model.User, messageEvent *model.MessageEvent) (bool, bool, string) {
	if messageEvent.TargetUserID != botUser.ID {
		msg := fmt.Sprintf("messageEvent is ignored because event is not for bot(id: %d) forUserID: %d", botUser.ID, messageEvent.TargetUserID)
		return false, false, msg
	}
	if messageEvent.IsShared {
		return false, false, "message event is ignored because it is shared event"
	}
	return IsTargetMessage(botUser, messageEvent.Message)
}

// IsTargetMessage checks if the given message is a target.
func IsTargetMessage(botUser *model.User, message *model.Message) (bool, bool, string) {
	if botUser.IsOwnMessage(message) {
		return false, false, "message is ignored because it is sent by bot"
	}

	if !message.IsRepliedTo(botUser) {
		reason := fmt.Sprintf("tweet is ignored because it is not sent to subscribe user(%d): receiver(%d)", botUser.ID, message.ReplyUserID)
		return false, false, reason
	}

	if message.ReferencedMessage != nil && botUser.IsOwnMessage(message.ReferencedMessage) {
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
