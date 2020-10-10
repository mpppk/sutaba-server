package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"
)

func ToTweet(anacondaTweet *anaconda.Tweet) *itwitter.Tweet {
	if anacondaTweet == nil {
		return nil
	}
	mediaList := getMediaList(anacondaTweet)
	var mediaURLs []string
	for _, media := range mediaList {
		mediaURLs = append(mediaURLs, media.Media_url_https)
	}

	tweet := &itwitter.Tweet{
		ID:                  anacondaTweet.Id,
		User:                *toUser(&anacondaTweet.User),
		Text:                anacondaTweet.Text,
		MediaURLs:           mediaURLs,
		InReplyToUserID:     anacondaTweet.InReplyToUserID,
		InReplyToStatusID:   anacondaTweet.InReplyToStatusID,
		InReplyToScreenName: anacondaTweet.InReplyToScreenName,
		RetweetedStatus:     ToTweet(anacondaTweet.RetweetedStatus),
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
