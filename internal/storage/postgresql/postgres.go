package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type FetchBody struct {
	Ticker     string  `json:"ticker"`
	Price      float64 `json:"price"`
	Difference float64 `json:"difference"`
}

func New(host string, port int, user, password, dbName string) (*Storage, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, fmt.Errorf("database connetction: %w", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddTicker(ticker string, price, diff float64) error {
	req := `insert into tickers (ticker, saved_date, price, difference)
		values($1, $2, $3, $4)
		on conflict (saved_date) do update
		set ticker = excluded.ticker
			saved_date = excluded.saved_date
			price = excluded.price
			difference = excluded.difference`
	date := time.Now().Format("2006-01-02")
	_, err := s.db.Exec(req, ticker, date, price, diff)
	return err
}

func (s *Storage) FetchParams(ticker, dateString string) (FetchBody, error) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, dateString)
	req := `select price, difference FROM tickers
		WHERE saved_date == $1`
	rows, err := s.db.Query(req, ticker, date)
	if err != nil {
		return FetchBody{}, err
	}

	var body FetchBody
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&body.Ticker, &body.Price, &body.Difference)
		if err != nil {
			return FetchBody{}, fmt.Errorf("scaning rows: %w", err)
		}
	}

	return body, nil
}
