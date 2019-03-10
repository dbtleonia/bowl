package bowl

import (
	"context"
	"google.golang.org/appengine/datastore"
)

// Outcome is stored under /Season,S/Outcome,B.
type Outcome struct {
	Winner string // teamID
}

func getOutcome(c context.Context, season, bowlID string) (*Outcome, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	outcomeKey := datastore.NewKey(c, "Outcome", bowlID, 0, seasonKey)
	var outcome Outcome
	if err := datastore.Get(c, outcomeKey, &outcome); err != nil {
		return nil, err
	}
	return &outcome, nil
}

func putOutcome(c context.Context, season, bowlID string, outcome *Outcome) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	outcomeKey := datastore.NewKey(c, "Outcome", bowlID, 0, seasonKey)
	_, err := datastore.Put(c, outcomeKey, outcome)
	return err
}

func deleteOutcome(c context.Context, season, bowlID string) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	outcomeKey := datastore.NewKey(c, "Outcome", bowlID, 0, seasonKey)
	return datastore.Delete(c, outcomeKey)
}

func readOutcomes(c context.Context, season string) (map[string]Outcome, error) {
	// TODO: Use cursor instead of GetAll.
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("Outcome").Ancestor(seasonKey)
	var outcomes []Outcome
	outcomeKeys, err := q.GetAll(c, &outcomes)
	if err != nil {
		return nil, err
	}
	outcomeMap := make(map[string]Outcome)
	for i, outcome := range outcomes {
		k := outcomeKeys[i]
		bowl := k.StringID()
		outcomeMap[bowl] = outcome
	}
	return outcomeMap, nil
}
