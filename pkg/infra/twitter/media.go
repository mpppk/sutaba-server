package twitter

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/domain/twitter"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/xerrors"
)

func DownloadEntityMedia(entityMedia *anaconda.EntityMedia, retryNum, retryInterval int) ([]byte, error) {
	mediaRawUrl := entityMedia.Media_url_https
	mediaUrl, err := url.Parse(mediaRawUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse media url(%s): %w", mediaRawUrl, err)
	}

	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return nil, xerrors.Errorf("invalid mediaUrl: %s", mediaRawUrl)
	}

	cnt := 0
	for {
		bytes, err := util.DownloadFile(mediaRawUrl)
		if err != nil {
			if cnt >= retryNum {
				return nil, xerrors.Errorf("failed to download file from %s (retired %d times): %w", mediaRawUrl, retryNum, err)
			}

			fmt.Println(xerrors.Errorf("failed to download file from %s: %w", mediaRawUrl, err))
			time.Sleep(time.Duration(retryInterval) * time.Second)
			cnt++
			continue
		}
		return bytes, nil
	}
}

func getMedia(tweet *anaconda.Tweet) (*anaconda.EntityMedia, bool) {
	entityMediaList := tweet.Entities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return nil, false
	}
	return &entityMediaList[0], true
}

func getMediaList(tweet *anaconda.Tweet) []anaconda.EntityMedia {
	entityMediaList := tweet.ExtendedEntities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return nil
	}
	return entityMediaList
}

func DownloadEntityMediaFromTweet(tweet *anaconda.Tweet, retryNum, retryInterval int) ([]byte, error) {
	media, ok := getMedia(tweet)
	if !ok {
		return nil, xerrors.Errorf("tweet(id: %d) has no media", tweet.Id)
	}
	return DownloadEntityMedia(media, retryNum, retryInterval)
}

func DownloadMediaFromTweet(tweet *twitter.Tweet, retryNum, retryInterval int) ([]byte, error) {
	if len(tweet.MediaURLs) == 0 {
		return nil, errors.New("tweet has no media")
	}

	mediaRawUrl := tweet.MediaURLs[0]
	util.LogPrintlnInOneLine("media URL:", mediaRawUrl)
	mediaUrl, err := url.Parse(mediaRawUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse media url(%s): %w", mediaRawUrl, err)
	}

	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return nil, xerrors.Errorf("invalid mediaUrl: %s", mediaRawUrl)
	}

	cnt := 0
	for {
		bytes, err := util.DownloadFile(mediaRawUrl)
		if err != nil {
			if cnt >= retryNum {
				return nil, xerrors.Errorf("failed to download file from %s (retired %d times): %w", mediaRawUrl, retryNum, err)
			}

			fmt.Println(xerrors.Errorf("failed to download file from %s: %w", mediaRawUrl, err))
			time.Sleep(time.Duration(retryInterval) * time.Second)
			cnt++
			continue
		}
		return bytes, nil
	}
}
