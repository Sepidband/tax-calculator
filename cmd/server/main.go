package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"tax-calculator/internal/api"
	"tax-calculator/internal/calculator"
	"tax-calculator/internal/client"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Get configuration
	apiURL := getEnv("TAX_API_URL", "http://localhost:5001")
	port := getEnv("PORT", "8080")

	// Initialize components
	taxClient := client.NewTaxAPIClient(apiURL)
	calc := calculator.New()
	handler := api.NewHandler(taxClient, calc, logger)

	// Setup routes
	router := gin.Default()

	// CORS middleware BEFORE your routes
	router.Use(CORSMiddleware())

	// Serve static files
	router.Static("/static", "./frontend")
	router.StaticFile("/", "./frontend/index.html")
	router.POST("/api/v1/calculate-tax", handler.CalculateTax)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	logger.Infof("Starting server on port %s", port)
	log.Fatal(router.Run(":" + port))
}

// CORS middleware function
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
