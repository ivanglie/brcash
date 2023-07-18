package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type Branches struct {
	Currency string   `json:"currency"`
	City     string   `json:"city"`
	Items    []Branch `json:"items"`
}

type Branch struct {
	Bank    string    `json:"bank"`
	Subway  string    `json:"subway"`
	Buy     float64   `json:"buy"`
	Sell    float64   `json:"sell"`
	Updated time.Time `json:"updated"`
}

type CurrencyCode string
type Region string

var (
	Moscow Region = "moskva"

	USD CurrencyCode = "840"
	EUR CurrencyCode = "978"
	AED CurrencyCode = "784"
	BYN CurrencyCode = "933"
	CAD CurrencyCode = "124"
	CHF CurrencyCode = "756"
	CNY CurrencyCode = "156"
	GBP CurrencyCode = "826"
	HKD CurrencyCode = "344"
	JPY CurrencyCode = "392"
	KRW CurrencyCode = "410"
	KTZ CurrencyCode = "398"
	TRY CurrencyCode = "949"

	CurrencyCodeMap = map[string]CurrencyCode{
		"USD": USD,
		"EUR": EUR,
		"AED": AED,
		"BYN": BYN,
		"CAD": CAD,
		"CHF": CHF,
		"CNY": CNY,
		"GBP": GBP,
		"HKD": HKD,
		"JPY": JPY,
		"KRW": KRW,
		"KTZ": KTZ,
		"TRY": TRY,
	}
)

// newBranch creates a new Branch instance.
func newBranch(bank, subway string, buy, sell float64, updated time.Time) Branch {
	return Branch{Bank: bank, Subway: subway, Buy: buy, Sell: sell, Updated: updated}
}

// String representation of cash currency exchange rates.
func (r *Branches) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}
