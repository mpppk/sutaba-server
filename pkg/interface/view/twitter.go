package view

type TwitterView interface {
	Show(text string) error
	ReplyToTweet(text string, toTweetIDStr string) error
}
