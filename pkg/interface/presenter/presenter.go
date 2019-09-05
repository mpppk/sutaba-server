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
	if err := r.view.ReplyToTweet(newText, message.ID); err != nil {
		return xerrors.Errorf("failed to reply message on Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyResultToMessageWithReference(toUser model.User, targetMessage, referredMessage *model.Message, result *domain.ClassifyResult) error {
	text := string(generateResultMessage(result))
	newText := toQuoteTweet(text, referredMessage.ID, referredMessage.User.Name)
	newText = toReplyText(newText, []model.UserName{toUser.Name})
	if err := r.view.ReplyToTweet(newText, targetMessage.ID); err != nil {
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

func toQuoteTweet(text string, quotedMessageID model.MessageID, quotedTweetUserScreenName model.UserName) string {
	return string(text) + " " + buildTweetUrl(
		quotedTweetUserScreenName,
		quotedMessageID,
	)
}

func buildTweetUrl(userName model.UserName, messageId model.MessageID) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d", userName, messageId)
}
