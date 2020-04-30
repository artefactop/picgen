package main

import (
	"net/http"

	"github.com/artefactop/picgen/internal/picgen"
)

func main() {

	http.HandleFunc("/", picgen.RootHandler)
	http.ListenAndServe(":8080", nil)
}
