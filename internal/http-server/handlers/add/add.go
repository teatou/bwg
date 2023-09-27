package add

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type TickerAdder interface {
	AddTicker(ticker string, price, diff float64) error
}

type ApiAddBody struct {
	Ticker string `json:"symbol"`
	// Date
	Price      float64 `json:"lastPrice"`
	Difference float64 `json:"priceChangePercent"`
}

type RequestAdd struct {
	Ticker string `json:"ticker"`
}

func New(tickerAdder TickerAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAdd

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		url := "https://api.binance.com/api/v3/ticker/24hr"
		resp, err := http.Get(url)
		if err != nil {
			render.JSON(w, r, fmt.Errorf("invalid request"))

			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			render.JSON(w, r, fmt.Errorf("invalid request"))

			return
		}

		var body []ApiAddBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		for _, b := range body {
			if b.Ticker == req.Ticker {
				err = tickerAdder.AddTicker(b.Ticker, b.Difference, b.Price)
				if err != nil {
					render.JSON(w, r, fmt.Errorf("invalid request"))

					return
				}

				render.JSON(w, r, Response{
					Status: "OK",
					Error:  "",
				})
			}
		}

		render.JSON(w, r, fmt.Errorf("invalid request"))
	}
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
