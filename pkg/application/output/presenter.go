package output

import (
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type MessagePresenter interface {
	Post(user model.TwitterUser, text string) error
	Reply(toUser model.TwitterUser, toMessageID, text string) error
	ReplyWithQuote(toUser model.TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName, text string) error
}
