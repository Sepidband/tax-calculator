// cmd/server/main.go
package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "tax-calculator/docs"
	"tax-calculator/internal/api"
	"tax-calculator/internal/calculator"
	"tax-calculator/internal/client"
)

// @title Tax Calculator API
// @version 1.0
// @description Canadian income tax calculator API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

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

	// CORS middleware BEFORE routes
	router.Use(CORSMiddleware())

	// Serve static files and API routes
	router.Static("/static", "./frontend")
	router.StaticFile("/", "./frontend/index.html")
	router.POST("/api/v1/calculate-tax", handler.CalculateTax)
	// Simple heath check
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "tax-calculator",
			"timestamp": time.Now().Unix(),
		})
	})

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server

	logger.Infof("Starting server on port %s", port)
	logger.Infof("Swagger docs available at: http://localhost:%s/swagger/index.html", port)
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
