package classifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mpppk/sutaba-server/pkg/util"

	"github.com/mpppk/sutaba-server/pkg/application/repository"

	"golang.org/x/xerrors"
)

type ImageClassifyServerRepository struct {
	host string
}

func NewImageClassifierServerRepository(host string) *ImageClassifyServerRepository {
	return &ImageClassifyServerRepository{
		host: host,
	}
}

type ImagePredictResponse struct {
	Pred       string `json:"pred"`
	Confidence string `json:"confidence"`
}

func (c *ImageClassifyServerRepository) Do(imageBytes []byte) (*repository.ClassifyResult, error) {
	body, contentType, err := util.GenerateMultipartFormBody(imageBytes)
	if err != nil {
		return nil, xerrors.Errorf("failed to create multipart form: %w", err)
	}
	url := c.host + "/predict"
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, xerrors.Errorf("failed to http post to predict endpoint(%s): %w", url, err)
	}

	var predict ImagePredictResponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &predict); err != nil {
		return nil, xerrors.Errorf("failed to decode response(%s) of classifier server: %w", string(bodyBytes), err)
	}

	if err = resp.Body.Close(); err != nil {
		return nil, xerrors.Errorf("failed to close predict response body: %w", err)
	}

	conf, err := strconv.ParseFloat(predict.Confidence, 32)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse confidence(%s) to float: %w", predict.Confidence)
	}
	result := &repository.ClassifyResult{
		Class:      predict.Pred,
		Confidence: conf,
	}
	return result, nil
}
