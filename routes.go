package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/pat"
	"net/http"
)

func main() {
	r := pat.New()
	r.Get("/log/{log_id}", WithDB(http.HandlerFunc(readLog)))
	r.Put("/log", WithDB(http.HandlerFunc(createLog)))
	r.Patch("/log/{log_id}", WithDB(http.HandlerFunc(updateLog)))
	r.Delete("/log/{log_id}", WithDB(http.HandlerFunc(deleteLog)))

	http.ListenAndServe(":9080", context.ClearHandler(r))
}
