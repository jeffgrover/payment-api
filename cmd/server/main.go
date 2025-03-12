package main

import (
	"os"

	"github.com/jeffgrover/payment-api/internal/api"
	"github.com/jeffgrover/payment-api/internal/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Connect to database
	database, err := db.New("payments.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	// Create API server
	apiConfig := api.Config{
		Title:       "Payments API",
		Version:     "1.0.0",
		Description: "A lightweight payment processing API built with Go, Huma, and SQLite, inspired by Stripe and Square but with a more focused feature set.",
	}
	server := api.New(database, apiConfig)

	// Start server
	addr := ":8080"
	log.Info().Msg("Starting Payments API server")
	if err := server.Start(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
		os.Exit(1)
	}
}
