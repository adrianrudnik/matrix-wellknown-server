package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type Config struct {
	Bind       string
	SchemaRoot string
}

func configure() *Config {
	config := &Config{}

	// Bind
	bind := os.Getenv("BIND_ADDR")
	if bind == "" {
		bind = ":8080"
	}

	config.Bind = bind

	log.Info().
		Str("given", bind).
		Msg("Webserver bind configured")

	// Schema root
	schemaRoot := os.Getenv("SCHEMA_ROOT")
	if schemaRoot == "" {
		schemaRoot = "/var/schema"
	}

	schemaRootResolved, err := filepath.Abs(schemaRoot)
	if err != nil {
		log.Panic().Err(err).Msg("Schema root path not found")
	}

	config.SchemaRoot = schemaRootResolved

	log.Info().
		Str("given", schemaRoot).
		Str("resolved", schemaRootResolved).
		Msg("Schema root configured")

	return config
}
