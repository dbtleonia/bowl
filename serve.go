package bowl

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
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

type IndexData struct {
	Bowls    map[string][]*Bowl                    // season -> bowls
	Players  map[string]Player                     // user -> player
	Picks    map[string]map[string]map[string]Pick // season -> user -> bowl -> pick
	Outcomes map[string]map[string]Outcome         // season -> bowl -> outcome
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFiles("templates/root.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, &IndexData{bowls, playerMap, pickMap, outcomeMap}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
