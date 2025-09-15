package tests

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "tax-calculator/internal/calculator"
    "tax-calculator/pkg/models"
)

func TestCalculateTax(t *testing.T) {
    calc := calculator.New()
    
    // Sample 2022 brackets (simplified)
    brackets := []models.TaxBracket{
        {Min: 0, Max: float64Ptr(50197), Rate: 0.15},
        {Min: 50197, Max: float64Ptr(100392), Rate: 0.205},
        {Min: 100392, Rate: 0.26}, // Top bracket
    }
    
    tests := []struct {
        name     string
        salary   float64
        expected float64
    }{
        {"Zero salary", 0, 0},
        {"Basic bracket", 50000, 7500},
        {"Two brackets", 100000, 17739.17},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calc.CalculateTax(tt.salary, brackets)
            assert.NoError(t, err)
            assert.InDelta(t, tt.expected, result.TotalTax, 0.01)
        })
    }
}

func float64Ptr(f float64) *float64 {
    return &f
}