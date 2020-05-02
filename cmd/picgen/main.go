package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/artefactop/picgen/internal/server"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", server.RootHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/{size}/{color}/{labelColor}", server.PathHandler).Methods("GET", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe(":3001", r))
}
