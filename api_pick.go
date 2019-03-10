package bowl

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func apiPickHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	season := vars["season"]
	bowlID := vars["bowl"]

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil || u.ID == "" {
		http.Error(w, "Must log in with Google account", http.StatusUnauthorized)
		return
	}
	if _, err := getPlayer(c, u.ID); err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.Error(w, "Must register", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	if time.Now().After(bowl.Kickoff) {
		http.Error(w, "Bowl has already started", http.StatusForbidden)
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
		if err := putPick(c, season, u.ID, bowlID, &Pick{Winner: teamID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		if err := deletePick(c, season, u.ID, bowlID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
