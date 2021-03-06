package itwitter

import (
	"strconv"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

type Tweet struct {
	ID                  int64
	User                model.User
	Text                string
	QuoteTweet          *Tweet
	InReplyToStatusID   int64
	InReplyToUserID     int64
	InReplyToScreenName string
	MediaURLs           []string
	RetweetedStatus     *Tweet
}

func (t *Tweet) GetIDStr() string {
	return strconv.FormatInt(t.ID, 10)
}

func (t *Tweet) IsReply() bool {
	return t.InReplyToUserID != 0
}

func (t *Tweet) HasQuoteTweet() bool {
	return t.QuoteTweet != nil
}

func (t *Tweet) GetFirstMediaURL() (string, bool) {
	if len(t.MediaURLs) == 0 {
		return "", false
	}
	return t.MediaURLs[0], true
}
