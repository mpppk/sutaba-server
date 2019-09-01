package usecase

import (
	"fmt"
	"strings"
	"time"

	twitter2 "github.com/mpppk/sutaba-server/pkg/domain/twitter"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/mpppk/sutaba-server/pkg/domain/message"

	"golang.org/x/xerrors"

	"github.com/mpppk/sutaba-server/pkg/infra/classifier"

	"github.com/mpppk/sutaba-server/pkg/infra/twitter"
)

type PostPredictTweetUseCaseConfig struct {
	SendUser          twitter2.TwitterUser
	TargetKeyword     string
	ClassifierClient  *classifier.Classifier
	ErrorTweetMessage string
	SorryTweetMessage string
	TwitterRepository twitter2.Repository
}

type PostPredictTweetUseCase struct {
	conf              *PostPredictTweetUseCaseConfig
	twitterRepository twitter2.Repository
}

func NewPostPredictTweetUsecase(conf *PostPredictTweetUseCaseConfig) *PostPredictTweetUseCase {
	return &PostPredictTweetUseCase{
		conf:              conf,
		twitterRepository: conf.TwitterRepository,
	}
}

func (p *PostPredictTweetUseCase) isTargetTweet(tweet *twitter2.Tweet) (bool, string) {
	if len(tweet.MediaURLs) == 0 {
		return false, "tweet is ignored because it has no media"
	}
	if !strings.Contains(tweet.Text, p.conf.TargetKeyword) {
		return false, "tweet is ignored because it has no keyword"
	}

	if tweet.User.ID == p.conf.SendUser.ID {
		return false, "tweet is ignored because it is sent by bot"
	}
	return true, ""
}

func (p *PostPredictTweetUseCase) ReplyToUser(tweet *twitter2.Tweet) (*twitter2.Tweet, string, error) {
	ok, reason := p.isTargetTweet(tweet)
	if ok {
		f := func() (*twitter2.Tweet, error) {
			tweetText, err := p.tweetToPredText(tweet)
			if err != nil {
				return nil, err
			}
			return p.twitterRepository.ReplyWithQuote(
				p.conf.SendUser,
				tweet.User,
				tweet.GetIDStr(),
				tweet.GetIDStr(),
				tweet.User.ScreenName,
				tweetText,
			)
		}
		postedTweet, err := f()
		if err != nil {
			errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())

			if _, err := p.twitterRepository.Post(p.conf.SendUser, errTweetText); err != nil {
				util.LogPrintlnInOneLine("failed to tweet error notify message", err)
			}

			if _, err := p.twitterRepository.Post(p.conf.SendUser, p.conf.SorryTweetMessage); err != nil {
				util.LogPrintlnInOneLine("failed to tweet error notify message", err)
			}
			return nil, "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
		}
		return postedTweet, "", nil
	}

	if !tweet.HasQuoteTweet() {
		return nil, reason, nil
	}

	// Check quote tweet
	ok, quoteReason := p.isTargetTweet(tweet.QuoteTweet)
	if !ok {
		return nil, reason + ", and " + quoteReason, nil
	}
	f := func() (*twitter2.Tweet, error) {
		tweetText, err := p.tweetToPredText(tweet.QuoteTweet)
		if err != nil {
			return nil, err
		}
		postedTweet, err := p.twitterRepository.ReplyWithQuote(
			p.conf.SendUser,
			tweet.User,
			tweet.GetIDStr(),
			tweet.QuoteTweet.GetIDStr(),
			tweet.QuoteTweet.User.ScreenName,
			tweetText,
		)

		if err != nil {
			return nil, xerrors.Errorf("failed to post tweet: %v", err)
		}
		return postedTweet, nil
	}
	postedTweet, err := f()
	if err != nil {
		errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
		if _, err := p.twitterRepository.Post(p.conf.SendUser, errTweetText); err != nil {
			util.LogPrintlnInOneLine("failed to tweet error notify message", err)
		}

		if _, err := p.twitterRepository.Post(p.conf.SendUser, p.conf.SorryTweetMessage); err != nil {
			util.LogPrintlnInOneLine("failed to tweet error notify message", err)
		}
		return nil, "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
	}
	return postedTweet, "", nil
}

func (p *PostPredictTweetUseCase) tweetToPredText(tweet *twitter2.Tweet) (string, error) {
	mediaBytes, err := twitter.DownloadMediaFromTweet(tweet, 3, 1)
	if err != nil {
		return "", err
	}

	predict, err := p.conf.ClassifierClient.Predict(mediaBytes)
	if err != nil {
		return "", xerrors.Errorf("failed to predict: %v", err)
	}

	tweetText, err := message.PredToText(predict)
	if err != nil {
		return "", xerrors.Errorf("failed to convert predict result to tweet text: %v", err)
	}
	return tweetText, err
}
