package bowl

import (
	"context"
	"google.golang.org/appengine/datastore"
)

// PlayerScore is stored under /Season,S/PlayerScore,U.
type PlayerScore struct {
	Correct    int
	Incorrect  int
	PctCorrect float64
	Wilson     float64
	Entropy    float64
	Maverick   float64
}

func writePlayerScore(c context.Context, season, userID string, score *PlayerScore) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	scoreKey := datastore.NewKey(c, "PlayerScore", userID, 0, seasonKey)
	_, err := datastore.Put(c, scoreKey, score)
	return err
}

func readPlayerScores(c context.Context, season string) (map[string]PlayerScore, error) {
	// TODO: Use cursor instead of GetAll
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").Ancestor(seasonKey)
	var playerScores []PlayerScore
	keys, err := q.GetAll(c, &playerScores)
	if err != nil {
		return nil, err
	}
	result := make(map[string]PlayerScore)
	for i, score := range playerScores {
		k := keys[i]
		userID := k.StringID()
		result[userID] = score
	}
	return result, nil
}

func readNumCorrectLeaders(c context.Context, season string) ([]string, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").
		Ancestor(seasonKey).
		Order("-Correct").
		KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, k := range keys {
		result = append(result, k.StringID())
	}
	return result, nil
}

func readPctCorrectLeaders(c context.Context, season string) ([]string, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").
		Ancestor(seasonKey).
		Order("-PctCorrect").
		KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, k := range keys {
		result = append(result, k.StringID())
	}
	return result, nil
}

func readWilsonLeaders(c context.Context, season string) ([]string, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").
		Ancestor(seasonKey).
		Order("-Wilson").
		KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, k := range keys {
		result = append(result, k.StringID())
	}
	return result, nil
}

func readEntropyLeaders(c context.Context, season string) ([]string, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").
		Ancestor(seasonKey).
		Order("Entropy").
		KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, k := range keys {
		result = append(result, k.StringID())
	}
	return result, nil
}

func readMaverickLeaders(c context.Context, season string) ([]string, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("PlayerScore").
		Ancestor(seasonKey).
		Order("-Maverick").
		KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, k := range keys {
		result = append(result, k.StringID())
	}
	return result, nil
}
