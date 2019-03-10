package bowl

import (
	"context"
	"google.golang.org/appengine/datastore"
)

// BowlScore is stored under /Season,S/BowlScore,B.
type BowlScore struct {
	Correct   int
	Incorrect int
}

func putBowlScore(c context.Context, season, bowlID string, score *BowlScore) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	scoreKey := datastore.NewKey(c, "BowlScore", bowlID, 0, seasonKey)
	_, err := datastore.Put(c, scoreKey, score)
	return err
}

func readBowlScores(c context.Context, season string) (map[string]BowlScore, error) {
	// TODO: Use cursor instead of GetAll
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("BowlScore").Ancestor(seasonKey)
	var bowlScores []BowlScore
	bowlScoreKeys, err := q.GetAll(c, &bowlScores)
	if err != nil {
		return nil, err
	}
	bowlScoreMap := make(map[string]BowlScore)
	for i, bowlScore := range bowlScores {
		k := bowlScoreKeys[i]
		bowl := k.StringID()
		bowlScoreMap[bowl] = bowlScore
	}
	return bowlScoreMap, nil
}

func (b BowlScore) PctIncorrect() float64 {
	return float64(b.Incorrect) / (float64(b.Correct) + float64(b.Incorrect))
}
