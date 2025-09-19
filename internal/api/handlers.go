// internal/api/handlers.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"tax-calculator/internal/calculator"
	"tax-calculator/internal/client"
)

type Handler struct {
	taxClient  *client.TaxAPIClient
	calculator *calculator.Calculator
	logger     *logrus.Logger
}

type CalculateTaxRequest struct {
	Salary float64 `json:"salary" binding:"required,min=0"`
	Year   int     `json:"year" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

func NewHandler(taxClient *client.TaxAPIClient, calc *calculator.Calculator, logger *logrus.Logger) *Handler {
	return &Handler{
		taxClient:  taxClient,
		calculator: calc,
		logger:     logger,
	}
}

// CalculateTax calculates income tax for a given salary and tax year
// @Summary Calculate income tax
// @Description Calculate total income tax, effective rate, and breakdown by tax brackets for Canadian tax system
// @Tags tax
// @Accept json
// @Produce json
// @Param request body CalculateTaxRequest true "Tax calculation request with annual salary and tax year"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Router /calculate-tax [post]
func (h *Handler) CalculateTax(c *gin.Context) {
	var req CalculateTaxRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"salary": req.Salary,
		"year":   req.Year,
	}).Info("Processing tax calculation request")

	// Get tax brackets
	brackets, err := h.taxClient.GetTaxBrackets(req.Year)
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch tax brackets")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Unable to fetch tax data"})
		return
	}

	// Calculate tax
	result, err := h.calculator.CalculateTax(req.Salary, brackets.TaxBrackets)
	if err != nil {
		h.logger.WithError(err).Error("Tax calculation failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Calculation failed"})
		return
	}

	h.logger.WithField("total_tax", result.TotalTax).Info("Tax calculation completed")
	c.JSON(http.StatusOK, result)
}
