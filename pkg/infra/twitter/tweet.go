package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

func ToTweet(anacondaTweet *anaconda.Tweet) *model.Tweet {
	mediaList := getMediaList(anacondaTweet)
	var mediaURLs []string
	for _, media := range mediaList {
		mediaURLs = append(mediaURLs, media.Media_url_https)
	}

	tweet := &model.Tweet{
		ID:        anacondaTweet.Id,
		User:      *toUser(&anacondaTweet.User),
		Text:      anacondaTweet.Text,
		MediaURLs: mediaURLs,
	}

	if anacondaTweet.InReplyToStatusID != 0 {
		tweet.InReplyToUserID = anacondaTweet.InReplyToUserID
		tweet.InReplyToScreenName = anacondaTweet.InReplyToScreenName
		tweet.InReplyToStatusID = anacondaTweet.InReplyToStatusID
	}

	if anacondaTweet.QuotedStatusID != 0 {
		tweet.QuoteTweet = ToTweet(anacondaTweet.QuotedStatus)
	}

	return tweet
}
