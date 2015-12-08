package bowl

import (
	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/schema"
)

type Player struct {
	Email    string
	Nickname string
}

type Pick struct {
	Winner string
}

type Outcome struct {
	Winner string
}

type BowlScore struct {
	Correct   int
	Incorrect int
}

type RootData struct {
	Bowls      map[string][]*Bowl                    // season -> bowls
	Players    map[string]Player                     // user -> player
	Picks      map[string]map[string]map[string]Pick // season -> user -> bowl -> pick
	Outcomes   map[string]map[string]Outcome         // season -> bowl -> outcome
	BowlScores map[string]map[string]BowlScore       // season -> bowl -> score
}

type UpdateBowlVars struct {
	Season string
	Bowl   string
}

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/updatebowl", updateBowlHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Player")
	var players []Player
	playerKeys, err := q.GetAll(c, &players)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playerMap := make(map[string]Player)
	for i, player := range players {
		k := playerKeys[i]
		user := k.StringID()
		playerMap[user] = player
	}

	q = datastore.NewQuery("Pick")
	var picks []Pick
	pickKeys, err := q.GetAll(c, &picks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pickMap := make(map[string]map[string]map[string]Pick)
	for i, pick := range picks {
		k := pickKeys[i]
		season := k.Parent().Parent().StringID()
		user := k.Parent().StringID()
		bowl := k.StringID()
		if _, present := pickMap[season]; !present {
			pickMap[season] = make(map[string]map[string]Pick)
		}
		if _, present := pickMap[season][user]; !present {
			pickMap[season][user] = make(map[string]Pick)
		}
		pickMap[season][user][bowl] = pick
	}

	q = datastore.NewQuery("Outcome")
	var outcomes []Outcome
	outcomeKeys, err := q.GetAll(c, &outcomes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	outcomeMap := make(map[string]map[string]Outcome)
	for i, outcome := range outcomes {
		k := outcomeKeys[i]
		season := k.Parent().StringID()
		bowl := k.StringID()
		if _, present := outcomeMap[season]; !present {
			outcomeMap[season] = make(map[string]Outcome)
		}
		outcomeMap[season][bowl] = outcome
	}

	q = datastore.NewQuery("BowlScore")
	var bowlScores []BowlScore
	bowlScoreKeys, err := q.GetAll(c, &bowlScores)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bowlScoreMap := make(map[string]map[string]BowlScore)
	for i, bowlScore := range bowlScores {
		k := bowlScoreKeys[i]
		season := k.Parent().StringID()
		bowl := k.StringID()
		if _, present := bowlScoreMap[season]; !present {
			bowlScoreMap[season] = make(map[string]BowlScore)
		}
		bowlScoreMap[season][bowl] = bowlScore
	}

	tmpl, err := template.ParseFiles("templates/root.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, &RootData{bowls, playerMap, pickMap, outcomeMap, bowlScoreMap}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateBowlHandler(w http.ResponseWriter, r *http.Request) {
	var vars UpdateBowlVars
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Reuse decoder.
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&vars, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := appengine.NewContext(r)
		k := datastore.NewKey(c, "Outcome", vars.Bowl, 0,
			datastore.NewKey(c, "Season", vars.Season, 0, nil))
		var outcome Outcome
		if err := datastore.Get(c, k, &outcome); err != nil {
			http.Error(w, fmt.Sprintf("No outcome for season=%q bowl=%q", vars.Season, vars.Bowl), http.StatusBadRequest)
			return
		}
		updateBowl.Call(c, vars.Season, vars.Bowl, &outcome)
	}

	tmpl, err := template.ParseFiles("templates/updatebowl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, vars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var updateBowl = delay.Func("key", func(c appengine.Context, season, bowl string, outcome *Outcome) {
	seasonKey := datastore.NewKey(c, "Season", season, 0, nil)

	// TODO: This query reads picks for all bowls.  Is there a way to filter?
	q := datastore.NewQuery("Pick").Ancestor(seasonKey)
	var picks []Pick
	pickKeys, err := q.GetAll(c, &picks)
	if err != nil {
		c.Warningf("%s", err)
		return
	}
	correct, incorrect := 0, 0
	for i, pick := range picks {
		k := pickKeys[i]
		if k.StringID() != bowl {
			continue
		}
		if pick.Winner == outcome.Winner {
			correct++
		} else {
			incorrect++
		}
	}

	scoreKey := datastore.NewKey(c, "BowlScore", bowl, 0, seasonKey)
	_, err = datastore.Put(c, scoreKey, &BowlScore{
		Correct:   correct,
		Incorrect: incorrect,
	})
	if err != nil {
		c.Warningf("%s", err)
		return
	}
})
