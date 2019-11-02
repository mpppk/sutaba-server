package ipresenter

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
)

type MessagePresenter interface {
	PostResult(result *domain.ClassifyResult, isDebug bool) error
	PostText(text string) error
	ReplyResultToMessage(message *model.Message, result *domain.ClassifyResult, isDebug bool) error
	ReplyToMessage(toMessage *model.Message, text string) error
	ReplyResultToMessageWithReference(targetMessage, referredMessage *model.Message, result *domain.ClassifyResult, isDebug bool) error
}
