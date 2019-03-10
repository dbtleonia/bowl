package bowl

import (
	"context"
	"google.golang.org/appengine/datastore"
)

// Player is stored under /Player,P.
type Player struct {
	Email    string
	Nickname string
}

func putPlayer(c context.Context, userID string, player *Player) error {
	playerKey := datastore.NewKey(c, "Player", userID, 0, nil)
	_, err := datastore.Put(c, playerKey, player)
	return err
}

func getPlayer(c context.Context, userID string) (*Player, error) {
	var player Player
	playerKey := datastore.NewKey(c, "Player", userID, 0, nil)
	if err := datastore.Get(c, playerKey, &player); err != nil {
		return nil, err
	}
	return &player, nil
}

func readPlayers(c context.Context) (map[string]Player, error) {
	// TODO: Use cursor instead of GetAll.
	q := datastore.NewQuery("Player")
	var players []Player
	playerKeys, err := q.GetAll(c, &players)
	if err != nil {
		return nil, err
	}
	playerMap := make(map[string]Player)
	for i, player := range players {
		k := playerKeys[i]
		user := k.StringID()
		playerMap[user] = player
	}
	return playerMap, nil
}
