package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
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

func GenerateMultipartFormBody(data []byte) (*bytes.Buffer, string, error) {
	dataBuffer := bytes.NewBuffer(data)
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	fw, err := mw.CreateFormFile("file", "image")

	// fwで作ったパートにファイルのデータを書き込む
	if _, err = io.Copy(fw, dataBuffer); err != nil {
		return nil, "", xerrors.Errorf("failed to copy image bytes to multipart form: %w", err)
	}

	// リクエストのContent-Typeヘッダに使う値を取得する（バウンダリを含む）
	contentType := mw.FormDataContentType()

	// 書き込みが終わったので最終のバウンダリを入れる
	if err = mw.Close(); err != nil {
		return nil, "", xerrors.Errorf("failed to close multipart writer: %w", err)
	}
	return body, contentType, nil
}
