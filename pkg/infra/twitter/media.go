package twitter

import (
	"github.com/ChimeraCoder/anaconda"
)

func getMediaList(tweet *anaconda.Tweet) []anaconda.EntityMedia {
	entityMediaList := tweet.ExtendedEntities.Media
	if entityMediaList == nil || len(entityMediaList) == 0 {
		return nil
	}
	return entityMediaList
}
