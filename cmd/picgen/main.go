package main

import (
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/artefactop/picgen/internal/server"
)

func main() {

	host, _ := os.Hostname()
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "picgen").
		Str("host", host).
		Logger()

	if err := profiler.Start(profiler.Config{
		Service:        "picgen",
		ServiceVersion: "0.0.1",
	}); err != nil {
		log.Warn().Msgf("Error starting profiling: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", server.RootHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/{size}/{color}/{labelColor}", server.PathHandler).Methods("GET", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal().Err(http.ListenAndServe(":3001", r))
}

func configureLogHandler(r *mux.Router) {
	logHandler := hlog.NewHandler(log.Logger)
	r.Use(logHandler)
	r.Use(hlog.AccessHandler((func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})))
	r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))
	r.Use(hlog.RequestIDHandler("req_id", "X-Request-Id"))
}
