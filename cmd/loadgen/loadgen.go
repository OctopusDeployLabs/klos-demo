package main

import (
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {
	// Config from env
	baseURL := os.Getenv("ENDPOINT")
	if baseURL == "" {
		baseURL = "http://worker"
	}

	rps := 1.0 // default to 1 request per second
	if rpsStr := os.Getenv("RPS"); rpsStr != "" {
		if parsed, err := strconv.ParseFloat(rpsStr, 32); err == nil {
			rps = parsed
		}
	}

	// Traffic split (percentage of requests to `/beverage` endpoint are for an iced coffee)
	coffeePct := 20 // default to 20%
	if pctStr := os.Getenv("COFFEE_PERCENT"); pctStr != "" {
		if parsed, err := strconv.Atoi(pctStr); err == nil && parsed >= 0 && parsed <= 100 {
			coffeePct = parsed
		}
	}

	logger.Info("Starting request sender", slog.String("endpoint", baseURL), slog.Int("ceo_percent", coffeePct))

	client := &http.Client{Timeout: 5 * time.Second}

	for {
		// Decide beverage choice
		beverage := "tea"
		route := "/beverage"
		if rand.Intn(100) < coffeePct {
			beverage = "coffee"
		}

		// Randomly choose between hot and cold drinks
		if rand.Intn(2) == 0 {
			route += "?kind=" + beverage + "&hot=false"
		} else {
			route += "?kind=" + beverage + "&hot=true"
		}
		url := baseURL + route

		// Send request
		sendRequest(client, url)

		time.Sleep(time.Duration(1 / rps * float64(time.Second))) // Pause between requests
	}
}

func sendRequest(client *http.Client, url string) {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)
	if err != nil {
		slog.Error("Request failed",
			slog.String("url", url),
			slog.String("error", err.Error()),
			slog.Duration("duration", duration),
		)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		slog.Info("Request unsuccessful",
			slog.String("url", url),
			slog.Int("status", resp.StatusCode),
			slog.Any("headers", resp.Header),
			slog.String("body", string(body)),
			slog.Duration("duration", duration),
		)
	} else {
		slog.Info("Request successful",
			slog.String("url", url),
			slog.Int("status", resp.StatusCode),
			slog.Any("headers", resp.Header),
			slog.String("body", string(body)),
			slog.Duration("duration", duration),
		)
	}
}
