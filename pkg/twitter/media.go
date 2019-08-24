package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

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
		bytes, err := DownloadFile(mediaRawUrl)
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

func DownloadEntityMediaFromTweet(tweet *anaconda.Tweet, retryNum, retryInterval int) ([]byte, error) {
	media, ok := getMedia(tweet)
	if !ok {
		return nil, xerrors.Errorf("tweet(id: %d) has no media", tweet.Id)
	}
	return DownloadEntityMedia(media, retryNum, retryInterval)
}

func DownloadFile(fileUrl string) (bytes []byte, err error) {
	response, err := http.Get(fileUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to request http get to %s: %w", fileUrl, err)
	}
	defer func() {
		cerr := response.Body.Close()
		if cerr == nil {
			return
		}
		err = xerrors.Errorf("failed to close http response body: %w", err)
	}()

	return ioutil.ReadAll(response.Body)
}
