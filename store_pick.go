package bowl

import (
	"context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Pick is stored under /Season,S/Player,U/Pick,B.
type Pick struct {
	Winner string // teamID
}

func putPick(c context.Context, season, userID, bowlID string, pick *Pick) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	playerKey := datastore.NewKey(c, "Player", userID, 0, seasonKey)
	bowlKey := datastore.NewKey(c, "Pick", bowlID, 0, playerKey)
	_, err := datastore.Put(c, bowlKey, pick)
	return err
}

func deletePick(c context.Context, season, userID, bowlID string) error {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	playerKey := datastore.NewKey(c, "Player", userID, 0, seasonKey)
	bowlKey := datastore.NewKey(c, "Pick", bowlID, 0, playerKey)
	return datastore.Delete(c, bowlKey)
}

func readPicksForBowl(c context.Context, season, bowlID string) ([]Pick, error) {
	// TODO: This query reads picks for all bowls.  Is there a way to filter?
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("Pick").Ancestor(seasonKey)
	var picks []Pick
	pickKeys, err := q.GetAll(c, &picks)
	if err != nil {
		return nil, err
	}
	var result []Pick
	for i, pick := range picks {
		k := pickKeys[i]
		if k.StringID() == bowlID {
			result = append(result, pick)
		}
	}
	return result, nil
}

func readPicksForUser(c context.Context, season, userID string) (map[string]Pick, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	playerKey := datastore.NewKey(c, "Player", userID, 0, seasonKey)
	q := datastore.NewQuery("Pick").Ancestor(playerKey)
	picks := make(map[string]Pick)
	for t := q.Run(c); ; {
		var pick Pick
		key, err := t.Next(&pick)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bowlID := key.StringID()
		picks[bowlID] = pick
	}
	return picks, nil
}

func readPicksForStartedBowls(c context.Context, season string, now time.Time) (map[string]map[string]Pick, error) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)
	q := datastore.NewQuery("Pick").Ancestor(seasonKey)
	var picks []Pick
	pickKeys, err := q.GetAll(c, &picks)
	if err != nil {
		return nil, err
	}
	pickMap := make(map[string]map[string]Pick)
	for i, pick := range picks {
		k := pickKeys[i]
		user := k.Parent().StringID()
		bowl := k.StringID()
		if now.Before(kickoff(season, bowl)) {
			continue
		}
		if _, present := pickMap[user]; !present {
			pickMap[user] = make(map[string]Pick)
		}
		pickMap[user][bowl] = pick
	}
	return pickMap, nil
}
