package bowl

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
)

type rootOutput struct {
	// Always present
	LoginURL  string
	LogoutURL string
	Season    string
	Bowls     []*Bowl

	// Present if logged in with Google account (nil otherwise)
	User *user.User

	// Present if logged in with Google account & registered (nil otherwise)
	Player *Player
	Picks  map[string]Pick
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	loginURL, _ := user.LoginURL(c, r.URL.String())
	logoutURL, _ := user.LogoutURL(c, r.URL.String())
	output := &rootOutput{
		LoginURL:  loginURL,
		LogoutURL: logoutURL,
		Season:    currentSeason,
		Bowls:     bowls[currentSeason],
	}
	u := user.Current(c)
	if u != nil && u.ID != "" {
		var err error
		output.User = u
		if output.Player, err = getPlayer(c, u.ID); err != nil {
			if err != datastore.ErrNoSuchEntity {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if output.Picks, err = readPicksForUser(c, currentSeason, u.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	tmpl, err := template.ParseFiles("templates/root.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
