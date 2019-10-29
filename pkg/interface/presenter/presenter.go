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

func (r *MessagePresenter) PostResult(result *domain.ClassifyResult) error {
	message, err := generateResultMessage(result)
	if err != nil {
		return xerrors.Errorf("failed to post result: %w", err)
	}
	return r.PostText(message)
}

func (r *MessagePresenter) PostText(text string) error {
	if err := r.view.Show(text); err != nil {
		return xerrors.Errorf("failed to post message to Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyToMessage(toMessage *model.Message, text string) error {
	newText := toReplyText(text, []model.UserName{toMessage.User.Name})
	if err := r.view.ReplyToTweet(newText, toMessage.ID); err != nil {
		return xerrors.Errorf("failed to reply message on Twitter: %w", err)
	}
	return nil
}

func (r *MessagePresenter) ReplyResultToMessage(toMessage *model.Message, result *domain.ClassifyResult) error {
	text, err := generateResultMessage(result)
	if err != nil {
		return xerrors.Errorf("failed to reply result to message: %w", err)
	}
	return r.ReplyToMessage(toMessage, text)
}

func (r *MessagePresenter) ReplyResultToMessageWithReference(targetMessage, referredMessage *model.Message, result *domain.ClassifyResult) error {
	text, err := generateResultMessage(result)
	if err != nil {
		return xerrors.Errorf("failed to post result: %w", err)
	}
	newText := toQuoteTweet(text, referredMessage.ID, referredMessage.User.Name)
	newText = toReplyText(newText, []model.UserName{targetMessage.User.Name})
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
