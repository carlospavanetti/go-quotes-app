package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome! Please hit the `/qod` API to get the quote of the day."))
}

func QuoteOfTheDayHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		date := currentTime.Format("2006-01-02")

		val, err := client.Get(date).Result()

		if err != redis.Nil {
			log.Println("Cache Hit for date ", date)
			w.Write([]byte(val))
			return
		}

		log.Println("Cache miss for date ", date)
		quoteResp, err := getQuoteFromAPI()
		if err != nil {
			w.Write([]byte("Sorry! We could not get the Quote of the Day. Please try again."))
			return
		}

		quote := quoteResp.Contents.Quotes[0].Quote
		client.Set(date, quote, 24*time.Hour)
		w.Write([]byte(quote))
	}
}

func getQuoteFromAPI() (*QuoteResponse, error) {
	API_URL := "http://quotes.rest/qod.json"
	resp, err := http.Get(API_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Quote API Returned: ", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("could not get quote from API")
	}

	quoteResp := &QuoteResponse{}
	json.NewDecoder(resp.Body).Decode(quoteResp)
	return quoteResp, nil
}
