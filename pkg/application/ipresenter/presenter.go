package ipresenter

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
)

type MessagePresenter interface {
	Post(user model.User, result *domain.ClassifyResult) error
	PostText(user model.User, text string) error
	Reply(toUser model.User, toMessageID string, result *domain.ClassifyResult) error
	ReplyWithQuote(toUser model.User, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName string, result *domain.ClassifyResult) error
}
