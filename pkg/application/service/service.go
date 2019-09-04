package service

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
	"github.com/mpppk/sutaba-server/pkg/util"
	"golang.org/x/xerrors"
)

// TODO: ここでいいのか考える
func DownloadMediaFromMessageMedia(media *model.MessageMedia, retryNum, retryInterval int) ([]byte, error) {
	mediaUrl, err := url.Parse(media.GetURL())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse media url(%s): %w", media, err)
	}

	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return nil, xerrors.Errorf("invalid mediaUrl: %s", media)
	}

	cnt := 0
	for {
		bytes, err := util.DownloadFile(media.GetURL())
		if err != nil {
			if cnt >= retryNum {
				return nil, xerrors.Errorf("failed to download file from %s (retired %d times): %w", media, retryNum, err)
			}

			fmt.Println(xerrors.Errorf("failed to download file from %s: %w", media, err))
			time.Sleep(time.Duration(retryInterval) * time.Second)
			cnt++
			continue
		}
		return bytes, nil
	}
}
