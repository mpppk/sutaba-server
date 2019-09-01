package twitter

import (
	"github.com/mpppk/sutaba-server/pkg/domain/twitter"
	"golang.org/x/xerrors"
)

type Repository struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
}

func NewRepository(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Repository {
	return &Repository{
		consumerKey:       consumerKey,
		consumerSecret:    consumerSecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
	}
}

func (r *Repository) Post(user twitter.TwitterUser, tweetText string) (*twitter.Tweet, error) {
	twitterUser := r.newUser(user.ID)
	postedTweet, err := twitterUser.Client.PostTweet(tweetText, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return ToTweet(&postedTweet), nil
}

func (r *Repository) Reply(fromUser, toUser twitter.TwitterUser, toTweetIDStr, tweetText string) (*twitter.Tweet, error) {
	fromTwitterUser := r.newUser(fromUser.ID)
	postedTweet, err := fromTwitterUser.PostReply(tweetText, toTweetIDStr, []string{toUser.ScreenName}) // FIXME toUserのIDがtwitterのIDとは限らない
	if err != nil {
		return nil, xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return postedTweet, nil
}

func (r *Repository) ReplyWithQuote(fromUser, toUser twitter.TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName, text string) (*twitter.Tweet, error) {
	fromTwitterUser := r.newUser(fromUser.ID)
	postedTweet, err := fromTwitterUser.PostReplyWithQuote(text, quotedTweetIDStr, quotedTweetUserScreenName, toTweetIDStr, []string{toUser.ScreenName}) // FIXME toUserのIDがtwitterのIDとは限らない
	if err != nil {
		return nil, xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return postedTweet, nil
}

func (r *Repository) newUser(id int64) *User {
	return NewUser(r.accessToken, r.accessTokenSecret, r.consumerKey, r.consumerSecret, id, "dummy keyword", true, ReplyWithQuote)
}
