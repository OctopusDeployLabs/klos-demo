package main

import (
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"klos-demo/pkg/handlers"
	"log/slog"
	"net/http"
	"os"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {

	loggingHandler := handlers.Logging(logger)
	r := gin.New()
	r.Use(loggingHandler)
	r.Use(gin.Recovery())
	r.GET("/", defaultHandler)
	r.GET("/ceo", ceoHandler)
	r.GET("/healthz", healthHandler)
	err := r.SetTrustedProxies(nil)
	if err != nil {
		logger.Error("error setting trusted proxies", "error", err)
	}
	err = r.Run(":8080")
	if err != nil {
		logger.Error("error starting server", "error", err)
	}
}

func healthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}
func defaultHandler(c *gin.Context) {
	c.Header("X-Correlation-Id", "default")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func ceoHandler(c *gin.Context) {
	c.Header("X-Correlation-Id", "ceo")
	enable, ok := os.LookupEnv("LetTheCeoDoThings")
	if !ok || enable != "true" {
		sloggin.AddCustomAttributes(c, slog.String("featureFlagRequired", "LetTheCeoDoThings"))
		c.JSON(500, gin.H{"error": "ceo is not allowed to do things"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
