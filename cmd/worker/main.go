package main

import (
	"github.com/gin-gonic/gin"
	"klos-demo/pkg/handlers"
	"log/slog"
	"net/http"
	"os"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

type Beverage struct {
	Kind string `form:"kind"`
	Hot  bool   `form:"hot"`
}

func main() {

	loggingHandler := handlers.Logging(logger)
	r := gin.New()
	r.Use(loggingHandler)
	r.Use(gin.Recovery())
	r.GET("/", defaultHandler)
	r.GET("/beverage", beverageHandler)
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

func beverageHandler(c *gin.Context) {

	beverageMachines := []string{"Espresso Machine", "Teapot", "Cold Brew"}

	var beverage Beverage
	err := c.ShouldBind(&beverage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	if beverage.Kind == "coffee" {
		if beverage.Hot {
			c.JSON(makeBeverage(beverage, beverageMachines[1]))
			return
		} else {
			// If the beverage is cold brew, use the Cold Brew
			c.JSON(makeBeverage(beverage, beverageMachines[2]))
			return
		}
	}
	if beverage.Kind == "tea" {
		c.JSON(makeBeverage(beverage, beverageMachines[1]))
		return
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Beverage not available"})
}

func makeBeverage(beverage Beverage, brewer string) (code int, obj any) {
	// Check that the brewer is valid
	if brewer == "Teapot" && beverage.Kind != "tea" {
		return http.StatusTeapot, gin.H{
			"error": "Cannot brew " + beverage.Kind + " in a teapot",
		}
	}

	// Get the temperature of the beverage
	var temperatureString string
	if beverage.Hot {
		temperatureString = "hot"
	} else {
		temperatureString = "cold"
	}

	return http.StatusOK, gin.H{
		"message": "Making a " + temperatureString + " " + beverage.Kind + " in " + brewer,
	}
}
