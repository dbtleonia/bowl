package bowl

import (
	"google.golang.org/appengine"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type outcomesOutput struct {
	Season   string
	Bowls    []*Bowl
	Outcomes map[string]Outcome
}

func outcomesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	season := vars["season"]

	c := appengine.NewContext(r)
	outcomes, err := readOutcomes(c, season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output := &outcomesOutput{
		Season:   season,
		Bowls:    bowls[season],
		Outcomes: outcomes,
	}
	tmpl, err := template.ParseFiles("templates/outcomes.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
