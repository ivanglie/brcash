package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type Branches struct {
	Currency Currency `json:"currency"`
	City     Region   `json:"city"`
	Items    []Branch `json:"items"`
}

type Branch struct {
	Bank    string    `json:"bank"`
	Subway  string    `json:"subway"`
	Buy     float64   `json:"buy"`
	Sell    float64   `json:"sell"`
	Updated time.Time `json:"updated"`
}

type Currency string
type Region string

var (
	USD    Currency = "USD"
	Moscow Region   = "moskva"
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
