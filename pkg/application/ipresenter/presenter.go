package ipresenter

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
)

type MessagePresenter interface {
	PostResult(user model.User, result *domain.ClassifyResult) error
	PostText(user model.User, text string) error
	ReplyResultToMessage(toUser model.User, message *model.Message, result *domain.ClassifyResult) error
	ReplyResultToMessageWithReference(toUser model.User, targetMessage, referredMessage *model.Message, result *domain.ClassifyResult) error
}
