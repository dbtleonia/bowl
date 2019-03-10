package bowl

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type leadersOutput struct {
	Season     string
	Bowls      []*Bowl
	Outcomes   map[string]Outcome   // bowlID -> outcome
	BowlScores map[string]BowlScore // bowlID -> score

	NumCorrect []string // userIDs
	PctCorrect []string // userIDs
	Wilson     []string // user IDs
	Entropy    []string // userIDs
	Maverick   []string // userIDs

	Players      map[string]Player          // userID -> player
	Picks        map[string]map[string]Pick // userID -> bowlID -> pick
	PlayerScores map[string]PlayerScore     // userID -> score
}

func leadersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	season := vars["season"]

	c := appengine.NewContext(r)
	outcomes, err := readOutcomes(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bowlScores, err := readBowlScores(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numCorrect, err := readNumCorrectLeaders(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pctCorrect, err := readPctCorrectLeaders(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wilson, err := readWilsonLeaders(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	entropy, err := readEntropyLeaders(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	maverick, err := readMaverickLeaders(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := readPlayers(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	picks, err := readPicksForStartedBowls(c, season, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playerScores, err := readPlayerScores(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := &leadersOutput{
		Season:     season,
		Bowls:      bowls[season],
		Outcomes:   outcomes,
		BowlScores: bowlScores,

		NumCorrect: numCorrect,
		PctCorrect: pctCorrect,
		Wilson:     wilson,
		Entropy:    entropy,
		Maverick:   maverick,

		Players:      players,
		Picks:        picks,
		PlayerScores: playerScores,
	}
	log.Errorf(ctx, "output: %v", output)
	tmpl, err := template.ParseFiles("templates/leaders.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
