package classifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/mpppk/sutaba-server/pkg/interface/itwitter"

	"github.com/mpppk/sutaba-server/pkg/domain/model"

	domain "github.com/mpppk/sutaba-server/pkg/domain/service"

	"github.com/mpppk/sutaba-server/pkg/util"

	"golang.org/x/xerrors"
)

type ImageClassifyServerService struct {
	host          string
	retryNum      int
	retryInterval int
	twitter       *itwitter.Twitter
}

func NewImageClassifierServerService(host string, retryNum, retryInterval int, twitter *itwitter.Twitter) *ImageClassifyServerService {
	return &ImageClassifyServerService{
		host:          host,
		retryNum:      retryNum,
		retryInterval: retryInterval,
		twitter:       twitter,
	}
}

type ImagePredictResponse struct {
	Pred       string `json:"pred"`
	Confidence string `json:"confidence"`
}

func (i *ImageClassifyServerService) Classify(message *model.Message) (*domain.ClassifyResult, error) {
	tweet, ok := i.twitter.RetrieveTweetFromMessage(message)
	if !ok {
		return nil, xerrors.Errorf("failed to retrieve tweet from message: %#v", message)
	}
	imageBytes, err := itwitter.DownloadMediaFromTweet(tweet, i.retryNum, i.retryInterval) // FIXME
	if err != nil {
		return nil, err
	}
	body, contentType, err := util.GenerateMultipartFormBody(imageBytes)
	if err != nil {
		return nil, xerrors.Errorf("failed to create multipart form: %w", err)
	}
	url := i.host + "/predict"
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
	result := &domain.ClassifyResult{
		Class:      predict.Pred,
		Confidence: conf,
	}
	return result, nil
}
