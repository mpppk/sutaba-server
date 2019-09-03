package output

import (
	"github.com/mpppk/sutaba-server/pkg/application/repository"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type MessagePresenter interface {
	Post(user model.TwitterUser, result *repository.ClassifyResult) error
	PostText(user model.TwitterUser, text string) error
	Reply(toUser model.TwitterUser, toMessageID string, result *repository.ClassifyResult) error
	ReplyWithQuote(toUser model.TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName string, result *repository.ClassifyResult) error
}
