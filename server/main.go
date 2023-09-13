package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeAPIExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ExchangeAPIResponse struct {
	USDBRL ExchangeAPIExchangeRate `json:"USDBRL"`
}

type ExchangeRate struct {
	gorm.Model
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

const (
	DOLLAR_EXCHANGE_API     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	DOLLAR_EXCHANGE_TIMEOUT = 200 * time.Millisecond
	DATABASE_TIMEOUT        = 10 * time.Millisecond
)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatalf("Fail to connect to database: %s", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/cotacao", NewExchangeRateHandler(db))

	http.ListenAndServe(":8080", mux)
}

type ExchangeRateHandler struct {
	db *gorm.DB
}

func NewExchangeRateHandler(db *gorm.DB) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		db: db,
	}
}

func (e ExchangeRateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := retrieveExchangeRate(ctx)
	if err != nil {
		log.Printf("Fail to retrieve Dollar exchange rate: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Fail to retrieve Dollar exchange rate",
		})
		return
	}

	exchangeRate := ExchangeRate{
		Code:       resp.USDBRL.Code,
		Codein:     resp.USDBRL.Codein,
		Name:       resp.USDBRL.Name,
		High:       resp.USDBRL.High,
		Low:        resp.USDBRL.Low,
		VarBid:     resp.USDBRL.VarBid,
		PctChange:  resp.USDBRL.PctChange,
		Bid:        resp.USDBRL.Bid,
		Ask:        resp.USDBRL.Ask,
		Timestamp:  resp.USDBRL.Timestamp,
		CreateDate: resp.USDBRL.CreateDate,
	}

	savedExchangeRate, err := saveExchangeRate(ctx, e.db, exchangeRate)
	if err != nil {
		log.Printf("Fail to save Dollar exchange rate: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Fail to save Dollar exchange rate",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedExchangeRate)
}

func retrieveExchangeRate(ctx context.Context) (*ExchangeAPIResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, DOLLAR_EXCHANGE_TIMEOUT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DOLLAR_EXCHANGE_API, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ExchangeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("server.sqlite"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&ExchangeRate{})

	return db, nil
}

func saveExchangeRate(
	ctx context.Context,
	db *gorm.DB,
	exchangeRate ExchangeRate,
) (*ExchangeRate, error) {
	ctx, cancelCtx := context.WithTimeout(ctx, DATABASE_TIMEOUT)
	defer cancelCtx()

	exchangeRateResponse := exchangeRate
	result := db.
		WithContext(ctx).
		Select("id", "created_at", "updated_at").
		Create(&exchangeRateResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	return &exchangeRateResponse, nil
}
