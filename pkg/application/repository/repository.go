package repository

import "github.com/mpppk/sutaba-server/pkg/domain/model"

type TwitterRepository interface {
	Post(user model.TwitterUser, text string) (*model.Tweet, error)
	Reply(fromUser, toUser model.TwitterUser, toTweetIDStr, text string) (*model.Tweet, error)
	ReplyWithQuote(fromUser, toUser model.TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName, text string) (*model.Tweet, error)
	DownloadMediaFromTweet(tweet *model.Tweet, retryNum, retryInterval int) ([]byte, error)
}

type ImageClassifierRepository interface {
	Do(image []byte) (*ClassifyResult, error)
}

type ClassifyResult struct {
	Class      string
	Confidence float64
}
