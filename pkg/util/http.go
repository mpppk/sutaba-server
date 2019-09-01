package util

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

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
