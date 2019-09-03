package presenter

import (
	"fmt"
	"strings"

	"github.com/mpppk/sutaba-server/pkg/application/repository"

	"github.com/mpppk/sutaba-server/pkg/application/service"
	"github.com/mpppk/sutaba-server/pkg/domain/model"
	"github.com/mpppk/sutaba-server/pkg/interface/view"
	"golang.org/x/xerrors"
)

type MessagePresenter struct {
	view view.MessageView
}

func (r *MessagePresenter) DownloadMediaFromTweet(tweet *model.Tweet, retryNum, retryInterval int) ([]byte, error) {
	return service.DownloadMediaFromTweet(tweet, retryNum, retryInterval)
}

func NewPresenter(view view.MessageView) *MessagePresenter {
	return &MessagePresenter{
		view: view,
	}
}

func (r *MessagePresenter) Post(user model.TwitterUser, result *repository.ClassifyResult) error {
	return r.PostText(user, generateResultMessage(result))
}

func (r *MessagePresenter) PostText(user model.TwitterUser, text string) error {
	if err := r.view.Show(text); err != nil {
		return xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) Reply(toUser model.TwitterUser, toTweetIDStr string, result *repository.ClassifyResult) error {
	text := generateResultMessage(result)
	newText := toReplyText(text, []string{toUser.ScreenName})
	if err := r.view.Reply(newText, &toUser); err != nil {
		return xerrors.Errorf("failed to reply message on Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyWithQuote(toUser model.TwitterUser, toTweetIDStr, quotedTweetIDStr, quotedTweetUserScreenName string, result *repository.ClassifyResult) error {
	text := generateResultMessage(result)
	newText := toQuoteTweet(text, quotedTweetIDStr, quotedTweetUserScreenName)
	newText = toReplyText(newText, []string{toUser.ScreenName})
	if err := r.view.Reply(newText, &toUser); err != nil {
		return xerrors.Errorf("failed to reply with quote message on Twitter: %w", err)
	}
	return nil
}

func toReplyText(text string, toScreenNames []string) string {
	var mentions []string
	for _, toScreenName := range toScreenNames {
		mentions = append(mentions, "@"+toScreenName)
	}
	return fmt.Sprintf("%s\n%s", strings.Join(mentions, " "), text)
}

func toQuoteTweet(text string, quotedTweetIDStr, quotedTweetUserScreenName string) string {
	return text + " " + buildTweetUrl(
		quotedTweetUserScreenName,
		quotedTweetIDStr,
	)
}

func buildTweetUrl(userName, id string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", userName, id)
}
