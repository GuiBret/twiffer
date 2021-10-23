package main

import (
	"net/http"
	"twitter-ripoff/users"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Use(jsonMiddleware)

	// dbmgmt.InitDB()
	r.HandleFunc("/users/{id:[0-9]+}", users.GetOneUser).Methods("GET")
	r.HandleFunc("/users", users.CreateUser).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":4000", nil)
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
