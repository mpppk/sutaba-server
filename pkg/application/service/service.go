package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
	"github.com/mpppk/sutaba-server/pkg/util"
	"golang.org/x/xerrors"
)

// TODO: ここでいいのか考える
func DownloadMediaFromTweet(tweet *model.Message, retryNum, retryInterval int) ([]byte, error) {
	if len(tweet.MediaURLs) == 0 {
		return nil, errors.New("tweet has no media")
	}

	mediaRawUrl := tweet.MediaURLs[0]
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
