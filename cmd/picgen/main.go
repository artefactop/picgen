package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/artefactop/picgen/internal/picgen"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{size}/{color}/{labelColor}", picgen.RootHandler).Methods("GET", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe(":3001", r))
}
