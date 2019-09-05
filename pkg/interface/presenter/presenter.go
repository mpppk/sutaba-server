package presenter

import (
	"fmt"
	"strings"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
	domain "github.com/mpppk/sutaba-server/pkg/domain/service"
	"github.com/mpppk/sutaba-server/pkg/interface/view"
	"golang.org/x/xerrors"
)

type MessagePresenter struct {
	view view.TwitterView
}

func NewPresenter(view view.TwitterView) *MessagePresenter {
	return &MessagePresenter{
		view: view,
	}
}

func (r *MessagePresenter) PostResult(user model.User, result *domain.ClassifyResult) error {
	return r.PostText(user, generateResultMessage(result))
}

func (r *MessagePresenter) PostText(user model.User, text string) error {
	if err := r.view.Show(text); err != nil {
		return xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyResultToMessage(toUser model.User, message *model.Message, result *domain.ClassifyResult) error {
	text := generateResultMessage(result)
	newText := toReplyText(text, []model.UserName{toUser.Name})
	if err := r.view.ReplyToTweet(newText, message.GetIDStr()); err != nil {
		return xerrors.Errorf("failed to reply message on Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyResultToMessageWithReference(toUser model.User, targetMessage, referredMessage *model.Message, result *domain.ClassifyResult) error {
	text := generateResultMessage(result)
	newText := toQuoteTweet(text, referredMessage.GetIDStr(), referredMessage.User.Name)
	newText = toReplyText(newText, []model.UserName{toUser.Name})
	if err := r.view.ReplyToTweet(newText, targetMessage.GetIDStr()); err != nil {
		return xerrors.Errorf("failed to reply with reference message on Twitter: %w", err)
	}
	return nil
}

func toReplyText(text string, toUserNames []model.UserName) string {
	var mentions []string
	for _, toScreenName := range toUserNames {
		mentions = append(mentions, "@"+string(toScreenName))
	}
	return fmt.Sprintf("%s\n%s", strings.Join(mentions, " "), text)
}

func toQuoteTweet(text string, quotedTweetIDStr string, quotedTweetUserScreenName model.UserName) string {
	return text + " " + buildTweetUrl(
		quotedTweetUserScreenName,
		quotedTweetIDStr,
	)
}

func buildTweetUrl(userName model.UserName, id string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", userName, id)
}
