package ipresenter

import (
	"github.com/mpppk/sutaba-server/pkg/application/repository"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type MessagePresenter interface {
	Post(user model.User, result *repository.ClassifyResult) error
	PostText(user model.User, text string) error
	Reply(toUser model.User, toMessageID string, result *repository.ClassifyResult) error
	ReplyWithQuote(toUser model.User, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName string, result *repository.ClassifyResult) error
}
