package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	client "github.com/mplewis/viteset-client-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var interval time.Duration

const defaultEndpoint = "https://api.viteset.com"
const defaultInterval = 15 * time.Second

// Key is a key that corresponds to a config blob.
type Key string

// mustEnv fetches an environment variable that is required to exist.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Str("var", key).Msg("Missing environment variable")
	}
	return val
}

// maybeEnv fetches an optional environment variable with a fallback value.
func maybeEnv(key string, dfault string) string {
	val := os.Getenv(key)
	if val == "" {
		return dfault
	}
	return val
}

// tty returns true if the current terminal is interactive.
func tty() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// init configures the application on startup.
func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Second
	if tty() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	interval = defaultInterval
	reqFresh := os.Getenv("FRESH")
	if reqFresh != "" {
		secs, err := strconv.ParseInt(reqFresh, 10, 0)
		if err != nil {
			log.Fatal().Err(err).Str("value", reqFresh).Msg("Invalid integer value for FRESH")
		}
		interval = time.Duration(secs) * time.Second
	}
}

// main starts the sidecar webserver.
func main() {
	blob := mustEnv("BLOB")
	endpoint := maybeEnv("ENDPOINT", defaultEndpoint)
	c := client.Client{Secret: mustEnv("SECRET"), Blob: blob, Host: endpoint, Interval: interval}
	ch, err := c.Subscribe()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to initialize client")
	}

	var val []byte = nil
	go func() {
		for {
			data := <-ch
			if data.Error != nil {
				log.Error().Err(err).Msg("Error while fetching blob value")
				return
			}
			val = data.Value
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(val)
	})

	host := maybeEnv("HOST", "0.0.0.0")
	port := maybeEnv("PORT", "8174")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Info().
		Str("address", addr).
		Str("endpoint", endpoint).
		Interface("blob", blob).
		Dur("interval", interval).
		Msg("Viteset Sidecar is ready")
	log.Fatal().Err(http.ListenAndServe(addr, nil)).Send()
}
