package main

import (
	"log"
	"net/http"

	"cloud.google.com/go/profiler"
	"github.com/gorilla/mux"

	"github.com/artefactop/picgen/internal/server"
)

func main() {

	if err := profiler.Start(profiler.Config{
		Service:        "picgen",
		ServiceVersion: "0.0.1",
	}); err != nil {
		log.Printf("Error starting profiling: %v\n", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", server.RootHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/{size}/{color}/{labelColor}", server.PathHandler).Methods("GET", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe(":3001", r))
}
