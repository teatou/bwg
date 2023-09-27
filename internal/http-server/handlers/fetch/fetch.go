package fetch

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/teatou/bwg/internal/storage/postgresql"
	"github.com/teatou/bwg/pkg/mylogger"
)

type Fetcher interface {
	FetchParams(ticker, dateString string) (postgresql.FetchBody, error)
}

func New(fetcher Fetcher, logger mylogger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticker := chi.URLParam(r, "ticker")
		if ticker == "" {
			logger.With("path", r.URL.Path).Errorf("error decoding query param 1")
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		dateString := chi.URLParam(r, "alias")
		if dateString == "" {
			logger.With("path", r.URL.Path).Errorf("error decoding query param 2")
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		body, err := fetcher.FetchParams(ticker, dateString)
		if err != nil {
			logger.With("path", r.URL.Path).Errorf("error fetching params: %v", err)
			render.JSON(w, r, fmt.Errorf("invalid request"))
			return
		}

		logger.With("path", r.URL.Path).Infof("successfully fetched params")
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
