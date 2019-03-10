package bowl

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/", rootHandler).
		Methods("GET")
	r.HandleFunc("/api/nickname", apiNicknameHandler).
		Methods("PUT", "DELETE")
	r.HandleFunc("/api/seasons/{season}/bowls/{bowl}/pick", apiPickHandler).
		Methods("PUT", "DELETE")
	r.HandleFunc("/leaders/{season}", leadersHandler).
		Methods("GET")

	r.HandleFunc("/admin/api/seasons/{season}/bowls/{bowl}/outcome", apiOutcomeHandler).
		Methods("PUT", "DELETE")
	r.HandleFunc("/admin/outcomes/{season}", outcomesHandler).
		Methods("GET")
	r.HandleFunc("/admin/updatebowl", updateBowlHandler).
		Methods("GET", "POST")
	r.HandleFunc("/admin/updateplayers", updatePlayersHandler).
		Methods("GET", "POST")

	http.Handle("/", r)
}
