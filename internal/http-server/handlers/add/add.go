package add

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/teatou/bwg/pkg/mylogger"
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

func New(tickerAdder TickerAdder, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAdd

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.With("path", r.URL.Path).Errorf("error decoding req body: %v", err)
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		url := "https://api.binance.com/api/v3/ticker/24hr"
		resp, err := http.Get(url)
		if err != nil {
			logger.With("path", r.URL.Path).Errorf("error req from binance api: %v", err)
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.With("path", r.URL.Path).Errorf("binance api status not ok: %v", err)
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		var body []ApiAddBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			logger.With("path", r.URL.Path).Errorf("error decoding binance body: %v", err)
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		for _, b := range body {
			if b.Ticker == req.Ticker {
				err = tickerAdder.AddTicker(b.Ticker, b.Difference, b.Price)
				if err != nil {
					logger.With("path", r.URL.Path).Errorf("error adding ticker: %v", err)
					render.JSON(w, r, fmt.Errorf("invalid request"))
					return
				}

				logger.With("path", r.URL.Path).Infof("success adding ticker: %v", err)
				render.JSON(w, r, Response{
					Status: "OK",
					Error:  "",
				})
			}
		}

		logger.With("path", r.URL.Path).Errorf("error?: %v", err)
		render.JSON(w, r, fmt.Errorf("invalid request"))
	}
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
