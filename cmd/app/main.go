package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

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
	api.Debug = opts.Dbg

	h := http.NewServeMux()
	h.HandleFunc("/", search)

	log.Info().Msg("Listening...")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Search...")

	client, err := api.NewClient()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	branches, err := client.Branches(r.URL.Query().Get("region"), r.URL.Query().Get("currency"))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	s := r.URL.Query().Get("sort")
	if len(s) == 0 {
		s = "by-sell"
	}

	b := branches.Items

	switch s {
	case "by-buy":
		log.Debug().Msg("Sort by buy")
		sort.Sort(sort.Reverse(api.ByBuySorter(b)))
		branches.Items = b
	case "by-sell":
		log.Debug().Msg("Sort by sell")
		sort.Sort(api.BySellSorter(b))
		branches.Items = b
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
