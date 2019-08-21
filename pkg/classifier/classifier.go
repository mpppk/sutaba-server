package classifier

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"golang.org/x/xerrors"
)

type Classifier struct {
	host string
}

func NewClassifier(host string) *Classifier {
	return &Classifier{
		host: host,
	}
}

type ImagePredictResponse struct {
	Pred       string `json:"pred"`
	Confidence string `json:"confidence"`
}

func (c *Classifier) Predict(imageBytes []byte) (*ImagePredictResponse, error) {
	body, contentType, err := generateMultipartFormBody(imageBytes)
	if err != nil {
		return nil, xerrors.Errorf("failed to create multipart form: %w", err)
	}
	url := c.host + "/predict"
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, xerrors.Errorf("failed to http post to predict endpoint(%s): %w", url, err)
	}

	var predict ImagePredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&predict); err != nil {
		remainBody, e := ioutil.ReadAll(resp.Body)
		remainBodyStr := ""
		if e == nil {
			remainBodyStr = string(remainBody)
		}
		return nil, xerrors.Errorf("failed to decode response(%s) of classifier server: %w", remainBodyStr, err)
	}
	if err = resp.Body.Close(); err != nil {
		return nil, xerrors.Errorf("failed to close predict response body: %w", err)
	}

	return &predict, nil
}

func generateMultipartFormBody(data []byte) (*bytes.Buffer, string, error) {
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
