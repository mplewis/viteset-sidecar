package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cache map[Key]*BlobData // Cached blob data store
var secret string           // The secret to use when requesting blobs
var endpoint string         // The Viteset API endpoint to use
var fresh time.Duration     // How long the app caches a blob
var onlyKey *string         // If set, app ignores the request path and always requests data for this blob

const defaultEndpoint = "https://api.viteset.com"
const defaultFresh = 15 * time.Second

// Key is a key that corresponds to a config blob.
type Key string

// BlobData is the data for a recently-fetched blob.
type BlobData struct {
	val   []byte
	stamp string
	at    time.Time
}

// fresh returns true if this blob has been fetched recently,
// along with the amount of time for which this blob will be fresh.
func (b *BlobData) fresh() (stillFresh bool, remain time.Duration) {
	expiry := b.at.Add(fresh)
	now := time.Now()
	if now.After(expiry) {
		return false, 0
	}
	return true, expiry.Sub(now)
}

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

// lookup fetches a blob from Viteset if the cached value is stale or missing.
func lookup(key Key) (found bool, val []byte, err error) {
	cached := cache[key]
	log := log.With().Str("key", string(key)).Bool("cached", cached != nil).Logger()

	if cached != nil {
		stillFresh, remain := cached.fresh()
		if stillFresh {
			log.Info().Dur("remain", remain).Msg("Blob is still fresh, not fetching")
			return true, cached.val, nil
		}
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", endpoint, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))
	if cached != nil {
		req.Header.Add("If-None-Match", cached.stamp)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}

	if resp.StatusCode == http.StatusNotModified && cached != nil {
		log.Info().Msg("Blob fetched and unchanged")
		cached.at = time.Now()
		return true, cached.val, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Warn().Msg("Blob not found")
		return false, nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("expected status 200 OK but got %d", resp.StatusCode)
	}

	stamp := resp.Header.Get("ETag")
	val, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}
	if cached == nil {
		cache[key] = &BlobData{val: val, stamp: stamp, at: time.Now()}
	} else {
		cached.val = val
		cached.stamp = stamp
		cached.at = time.Now()
	}
	log.Info().Msg("Blob fetched and updated")
	return true, val, nil
}

// tty returns true if the current terminal is interactive.
func tty() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// init configures the application on startup.
func init() {
	cache = map[Key]*BlobData{}

	secret = mustEnv("SECRET")
	endpoint = maybeEnv("ENDPOINT", defaultEndpoint)
	reqKey := os.Getenv("BLOB")
	if reqKey != "" {
		onlyKey = &reqKey
	}
	reqFresh := os.Getenv("FRESH")
	if reqFresh != "" {
		secs, err := strconv.ParseInt(reqFresh, 10, 0)
		if err != nil {
			log.Fatal().Err(err).Str("value", reqFresh).Msg("Invalid integer value for FRESH")
		}
		fresh = time.Duration(secs) * time.Second
	} else {
		fresh = defaultFresh
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Second
	if tty() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
}

// main starts the sidecar webserver.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]
		if onlyKey != nil {
			key = *onlyKey
		}

		found, val, err := lookup(Key(key))
		if err != nil {
			log.Error().Str("key", key).Err(err).Msg("Error fetching blob")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(val)
	})

	host := maybeEnv("HOST", "0.0.0.0")
	port := maybeEnv("PORT", "80")
	addr := fmt.Sprintf("%s:%s", host, port)
	blob := "<specified by requester>"
	if onlyKey != nil {
		blob = *onlyKey
	}
	log.Info().
		Str("address", addr).
		Str("endpoint", endpoint).
		Interface("blob", blob).
		Dur("fresh_secs", fresh).
		Msg("Viteset Sidecar is ready")
	log.Fatal().Err(http.ListenAndServe(addr, nil)).Send()
}
