package twitter

type Repository interface {
	Post(user TwitterUser, text string) (*Tweet, error)
	Reply(fromUser, toUser TwitterUser, toTweetIDStr, text string) (*Tweet, error)
	ReplyWithQuote(fromUser, toUser TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName, text string) (*Tweet, error)
}
