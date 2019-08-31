package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/domain/message"

	"golang.org/x/xerrors"

	"github.com/mpppk/sutaba-server/pkg/infra/classifier"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/sutaba-server/pkg/infra/twitter"
)

type PostPredictTweetUseCaseConfig struct {
	SendUser          *twitter.User
	ClassifierClient  *classifier.Classifier
	ErrorTweetMessage string
	SorryTweetMessage string
}

type PostPredictTweetUseCase struct {
	conf *PostPredictTweetUseCaseConfig
}

func NewPostPredictTweetUsecase(conf *PostPredictTweetUseCaseConfig) *PostPredictTweetUseCase {
	return &PostPredictTweetUseCase{
		conf: conf,
	}
}

func (p *PostPredictTweetUseCase) isTargetTweet(tweet *anaconda.Tweet) (bool, string) {
	entityMediaList := tweet.Entities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return false, "tweet is ignored because it has no media"
	}

	if !strings.Contains(tweet.Text, p.conf.SendUser.TargetKeyword) {
		return false, "tweet is ignored because it has no keyword"
	}

	if tweet.User.Id == p.conf.SendUser.ID {
		return false, "tweet is ignored because it is sent by bot"
	}
	return true, ""
}

func (p *PostPredictTweetUseCase) ReplyToUser(tweet *anaconda.Tweet) (*anaconda.Tweet, string, error) {
	ok, reason := p.isTargetTweet(tweet)
	if ok {
		postedTweet, err := p.postPredictTweet(tweet, "")
		if err != nil {
			errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
			p.conf.SendUser.PostErrorTweet(errTweetText, p.conf.SorryTweetMessage, tweet.IdStr, tweet.User.ScreenName)
			return nil, "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
		}
		return postedTweet, "", nil
	}

	if tweet.QuotedStatus == nil {
		return nil, reason, nil
	}

	// Check quote tweet
	ok, quoteReason := p.isTargetTweet(tweet.QuotedStatus)
	if !ok {
		return nil, reason + ", and " + quoteReason, nil
	}
	f := func() (*anaconda.Tweet, error) {
		tweetText, err := p.tweetToPredText(tweet.QuotedStatus)
		if err != nil {
			return nil, err
		}
		postedTweet, err := p.conf.SendUser.PostReplyWithQuote(tweetText, tweet.QuotedStatus, tweet.IdStr, []string{tweet.User.ScreenName})
		if err != nil {
			return nil, xerrors.Errorf("failed to post tweet: %v", err)
		}
		return &postedTweet, nil
	}
	postedTweet, err := f()
	if err != nil {
		errTweetText := p.conf.ErrorTweetMessage + fmt.Sprintf(" %v", time.Now())
		p.conf.SendUser.PostErrorTweet(errTweetText, p.conf.SorryTweetMessage, tweet.IdStr, tweet.User.ScreenName)
		return nil, "", xerrors.Errorf("error occurred in JudgeAndPostPredictTweetUseCase: %w", err)
	}
	return postedTweet, "", nil
}

func (p *PostPredictTweetUseCase) postPredictTweet(tweet *anaconda.Tweet, tweetTextPrefix string) (*anaconda.Tweet, error) {
	tweetText, err := p.tweetToPredText(tweet)
	if err != nil {
		return nil, err
	}
	postedTweet, err := p.conf.SendUser.PostByTweetType(tweetTextPrefix+tweetText, tweet)
	if err != nil {
		return nil, xerrors.Errorf("failed to post tweet: %v", err)
	}
	return &postedTweet, nil
}

func (p *PostPredictTweetUseCase) tweetToPredText(tweet *anaconda.Tweet) (string, error) {
	mediaBytes, err := twitter.DownloadEntityMediaFromTweet(tweet, 3, 1)
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
