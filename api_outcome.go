package bowl

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"

	"github.com/gorilla/mux"
)

func apiOutcomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	season := vars["season"]
	bowlID := vars["bowl"]

	c := appengine.NewContext(r)
	if !user.IsAdmin(c) {
		http.Error(w, "Must be an administrator", http.StatusUnauthorized)
		return
	}
	if _, ok := bowls[season]; !ok {
		http.Error(w, fmt.Sprintf("Invalid season: %q", season), http.StatusBadRequest)
		return
	}
	bowl, ok := bowlByID(bowls[season], bowlID)
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid bowl: %q", bowlID), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "PUT":
		var teamID string
		if err := json.NewDecoder(r.Body).Decode(&teamID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if teamID != bowl.Team1ID && teamID != bowl.Team2ID {
			http.Error(w, "Invalid team", http.StatusBadRequest)
			return
		}
		if err := putOutcome(c, season, bowlID, &Outcome{Winner: teamID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		if err := deleteOutcome(c, season, bowlID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
