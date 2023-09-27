package fetch

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/teatou/bwg/internal/storage/postgresql"
)

type Fetcher interface {
	FetchParams(ticker, dateString string) (postgresql.FetchBody, error)
}

func New(fetcher Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticker := chi.URLParam(r, "ticker")
		if ticker == "" {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		dateString := chi.URLParam(r, "alias")
		if dateString == "" {
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		body, err := fetcher.FetchParams(ticker, dateString)
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
