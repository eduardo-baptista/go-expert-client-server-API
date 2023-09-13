package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ExchangeRate struct {
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
	EXCHANGE_RATE_URL     = "http://localhost:8080/cotacao"
	EXCHANGE_RATE_TIMEOUT = 300 * time.Millisecond
	EXCHANGE_RATE_FILE    = "cotacao.txt"
)

func main() {
	exchangeRate, err := retrieveDollarExchangeRate(context.Background())
	if err != nil {
		log.Fatalf("could not retrieve dollar exchange rate: %v", err)
	}
	log.Printf("bid: %s", exchangeRate.Bid)

	err = saveExchangeRateOnFile(exchangeRate.Bid)
	if err != nil {
		log.Fatalf("could not save dollar exchange rate on file: %v", err)
	}
}

func retrieveDollarExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(ctx, EXCHANGE_RATE_TIMEOUT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, EXCHANGE_RATE_URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var exchangeRate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRate); err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func saveExchangeRateOnFile(exchangeRate string) error {
	file, err := os.Create(EXCHANGE_RATE_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: {%s}", exchangeRate))
	return err
}
