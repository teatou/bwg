package fetch

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/teatou/bwg/internal/storage/postgresql"
)

type Fetcher interface {
	FetchParams(ticker, dateString string) (postgresql.FetchBody, error)
}

type RequestFetch struct {
	Ticker     string `json:"ticker"`
	DateString string `json:"date"`
}

func New(fetcher Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestFetch

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		body, err := fetcher.FetchParams(req.Ticker, req.DateString)
		if err != nil {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		render.JSON(w, r, ResponseOK{
			Status: "OK",
			Body:   body,
		})
	}
}

type ResponseOK struct {
	Status string               `json:"status"`
	Body   postgresql.FetchBody `json:"body"`
}
