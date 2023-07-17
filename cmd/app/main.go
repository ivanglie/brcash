package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ivanglie/brcash/internal/api"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	opts struct {
		Dbg bool `long:"dbg" env:"DEBUG" description:"use debug"`
	}

	version = "unknown"
)

func main() {
	fmt.Printf("brcash %s\n", version)

	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] brcash error: %v", err)
		}
		os.Exit(2)
	}

	setupLog(opts.Dbg)

	h := http.NewServeMux()
	h.HandleFunc("/", search)

	log.Info().Msg("Listening...")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Search...")

	region := r.URL.Query().Get("region")
	if len(region) == 0 {
		region = string(api.Moscow)
	}

	log.Info().Msgf("region is %s", region)

	client, err := api.NewClient()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	branches, err := client.Branches(api.USD, api.Region(region))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(branches)
}

func setupLog(dbg bool) {
	if dbg {
		log.Level(zerolog.DebugLevel)
		return
	}

	log.Level(zerolog.InfoLevel)
}
