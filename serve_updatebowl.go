package bowl

import (
	"context"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"

	"github.com/gorilla/schema"
)

type updateBowlInput struct {
	Season string
	Bowl   string // TODO: BowlID
}

func updateBowlHandler(w http.ResponseWriter, r *http.Request) {
	var input updateBowlInput
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Reuse decoder.
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&input, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := appengine.NewContext(r)
		outcome, err := getOutcome(c, input.Season, input.Bowl)
		if err != nil {
			http.Error(w, fmt.Sprintf("No outcome for season=%q bowl=%q", input.Season, input.Bowl), http.StatusBadRequest)
			return
		}
		updateBowl.Call(c, input.Season, input.Bowl, outcome)
	}

	tmpl, err := template.ParseFiles("templates/updatebowl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var updateBowl = delay.Func("key", func(c context.Context, season, bowlID string, outcome *Outcome) {
	picks, err := readPicksForBowl(c, season, bowlID)
	if err != nil {
		log.Warningf(c, "%s", err)
		return
	}
	correct, incorrect := 0, 0
	for _, pick := range picks {
		if pick.Winner == outcome.Winner {
			correct++
		} else {
			incorrect++
		}
	}
	err = putBowlScore(c, season, bowlID, &BowlScore{
		Correct:   correct,
		Incorrect: incorrect,
	})
	if err != nil {
		log.Warningf(c, "%s", err)
		return
	}
})
