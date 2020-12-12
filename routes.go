package main

import (
	"fmt"
	cyphar "github.com/cyphar/filepath-securejoin"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// https://matrix.org/docs/spec/client_server/latest#get-well-known-matrix-client
// https://matrix.org/docs/spec/server_server/latest#get-well-known-matrix-server

func SchemaLookupRoute(w http.ResponseWriter, r *http.Request) {
	// Extract the requested schema taret
	target := chi.URLParam(r, "target")

	if target != "server" && target != "client" {
		log.Info().Str("target", target).Msg("Invalid schema target requested")
		w.WriteHeader(400)
	}

	config := r.Context().Value("Config").(*Config)

	// Extract host name only, remove possible port
	host := strings.Split(r.Host, ":")[0]

	// Create secured file path to possible schema file
	schema, err := cyphar.SecureJoin(config.SchemaRoot, host)
	schema = fmt.Sprintf("%s.%s.json", schema, target)

	// Detect invalid paths and exit
	if err != nil {
		log.Warn().Err(err).Str("schema", schema).Str("host", host).Msg("Invalid request path detected")
		w.WriteHeader(400)
	}

	// Detect undefined schemas
	if _, err := os.Stat(schema); os.IsNotExist(err) {
		w.WriteHeader(404)
		return
	}

	render.SetContentType(render.ContentTypeJSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	json, err := ioutil.ReadFile(schema)

	if err != nil {
		log.Error().Err(err).Str("schema", schema).Msg("Failed to read the resolved schema file")
	}

	_, err = w.Write(json)

	if err != nil {
		log.Error().Err(err).Str("schema", schema).Msg("Failed to write the resolved schema file to http writer")
	}
}
