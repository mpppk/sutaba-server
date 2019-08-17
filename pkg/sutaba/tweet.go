package sutaba

import (
	"fmt"
	"strconv"

	"golang.org/x/xerrors"
)

func PredToText(predict *ImagePredictResponse) (string, error) {
	confidence, err := strconv.ParseFloat(predict.Confidence, 32)
	if err != nil {
		return "", xerrors.Errorf("failed to parse confidence(%s) to float: %w", predict.Confidence)
	}
	predStr := ""
	switch predict.Pred {
	case "sutaba":
		if confidence > 0.8 {
			predStr = "間違いなくスタバ"
		} else if confidence > 0.5 {
			predStr = "スタバ"
		} else {
			predStr = "たぶんスタバ"
		}
	case "ramen":
		if confidence > 0.8 {
			predStr = "どう見てもラーメン"
		} else if confidence > 0.5 {
			predStr = "ラーメン"
		} else {
			predStr = "ラーメン...?"
		}
	case "other":
		if confidence > 0.8 {
			predStr = "スタバではない"
		} else if confidence > 0.5 {
			predStr = "スタバとは言えない"
		} else {
			predStr = "なにこれ...スタバではない気がする"
		}
	}

	return fmt.Sprintf("判定:%s\n確信度:%.2f", predStr, confidence*100) + "%", nil
}
