package bowl

import (
	"context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"html/template"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/schema"
)

type updatePlayersInput struct {
	Season string
}

func updatePlayersHandler(w http.ResponseWriter, r *http.Request) {
	var input updatePlayersInput
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
		updatePlayers.Call(c, input.Season)
	}

	tmpl, err := template.ParseFiles("templates/updateplayers.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var updatePlayers = delay.Func("key", func(c context.Context, season string) {
	outcomes, err := readOutcomes(c, season)
	if err != nil {
		log.Warningf(c, "%s, err")
		return
	}
	bowlScores, err := readBowlScores(c, season)
	if err != nil {
		log.Warningf(c, "%s", err)
		return
	}
	// TODO: Is this right?  Is "started bowls" the right set to read?
	picks, err := readPicksForStartedBowls(c, season, time.Now())
	if err != nil {
		log.Warningf(c, "%s", err)
		return
	}
	for userID, userPicks := range picks {
		correct, incorrect := 0, 0
		maverick := 0.0
		for bowlID, pick := range userPicks {
			if bowlScore, ok := bowlScores[bowlID]; ok {
				if pick.Winner == outcomes[bowlID].Winner {
					correct++
					if bowlScore.Correct+bowlScore.Incorrect > 0 {
						maverick += float64(bowlScore.Incorrect) / float64(bowlScore.Correct+bowlScore.Incorrect)
					}
				} else {
					incorrect++
				}
			}
		}
		pctCorrect := float64(correct) / float64(correct+incorrect)
		err := writePlayerScore(c, season, userID, &PlayerScore{
			Correct:    correct,
			Incorrect:  incorrect,
			PctCorrect: pctCorrect,
			Wilson:     wilson(correct, correct+incorrect),
			Entropy:    binaryEntropy(pctCorrect),
			Maverick:   maverick,
		})
		log.Warningf(c, "err =====> %v", err)
		if err != nil {
			log.Warningf(c, "%s", err)
			return
		}
	}
})

// http://www.evanmiller.org/how-not-to-sort-by-average-rating.html
func wilson(positive, total int) float64 {
	if total == 0 {
		return 0
	}
	const z = float64(1.96)
	pos, n := float64(positive), float64(total)
	phat := pos / n
	return (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
}

func binaryEntropy(p float64) float64 {
	h := 0.0
	if p > 0.0 {
		h += -p * math.Log2(p)
	}
	if p < 1.0 {
		h += -(1 - p) * math.Log2(1-p)
	}
	return h
}
