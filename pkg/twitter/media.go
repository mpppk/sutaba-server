package twitter

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/xerrors"
)

func DownloadEntityMedia(entityMedia *anaconda.EntityMedia) ([]byte, error) {
	mediaRawUrl := entityMedia.Media_url_https
	mediaUrl, err := url.Parse(mediaRawUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse media url(%s): %w", mediaRawUrl, err)
	}

	mediaUrlPaths := strings.Split(mediaUrl.Path, "/")
	if len(mediaUrlPaths) == 0 {
		return nil, xerrors.Errorf("invalid mediaUrl: %s", mediaRawUrl)
	}

	bytes, err := DownloadFile(mediaRawUrl)
	if err != nil {
		return nil, xerrors.Errorf("failed to download file to %s")
	}
	return bytes, nil
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
