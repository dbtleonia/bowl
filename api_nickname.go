package bowl

import (
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"
)

func apiNicknameHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil || u.ID == "" {
		http.Error(w, "Must log in with Google account", http.StatusUnauthorized)
		return
	}
	var nickname string
	if err := json.NewDecoder(r.Body).Decode(&nickname); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := putPlayer(c, u.ID, &Player{
		Email:    u.Email,
		Nickname: nickname,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
