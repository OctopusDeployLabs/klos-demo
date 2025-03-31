package main

import (
	"github.com/gin-gonic/gin"
	"klos-demo/pkg/handlers"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {
	go startFillingCache(300)

	loggingHandler := handlers.LoggingWithNoHealth(logger)
	r := gin.New()
	r.Use(loggingHandler)
	r.GET("/healthz", healthHandler)
	err := r.Run(":8080")
	if err != nil {
		logger.Error("error starting server", "error", err)
	}

}

func startFillingCache(sizeInMiB int) {
	sizeInBytes := sizeInMiB * 1024 * 1024

	logger := logger.With("sizeInBytes", sizeInBytes)

	logger.Info("Allocating memory")

	data := make([]byte, sizeInBytes)

	// Fill with pseudo-random data
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	logger.Info("Finished allocating memory")
}

func healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}
