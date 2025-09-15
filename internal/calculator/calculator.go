package calculator

import (
    "tax-calculator/pkg/models"
)

type Calculator struct{}

func New() *Calculator {
    return &Calculator{}
}

func (c *Calculator) CalculateTax(salary float64, brackets []models.TaxBracket) (*models.TaxCalculationResult, error) {
    if salary <= 0 {
        return &models.TaxCalculationResult{
            TotalTax:         0,
            EffectiveRate:    0,
            BracketBreakdown: []models.BracketTaxDetail{},
        }, nil
    }
    
    var totalTax float64
    var breakdown []models.BracketTaxDetail
    remainingSalary := salary
    
    for _, bracket := range brackets {
        if remainingSalary <= 0 {
            break
        }
        
        // Calculate taxable amount in this bracket
        var taxableInBracket float64
        
        if bracket.Max == nil {
            // Top bracket (unlimited)
            taxableInBracket = remainingSalary
        } else {
            bracketSize := *bracket.Max - bracket.Min
            taxableInBracket = min(remainingSalary, bracketSize)
        }
        
        if taxableInBracket > 0 {
            taxInBracket := taxableInBracket * bracket.Rate
            totalTax += taxInBracket
            
            breakdown = append(breakdown, models.BracketTaxDetail{
                Bracket: bracket,
                TaxOwed: taxInBracket,
            })
            
            remainingSalary -= taxableInBracket
        }
    }
    
    effectiveRate := 0.0
    if salary > 0 {
        effectiveRate = totalTax / salary
    }
    
    return &models.TaxCalculationResult{
        TotalTax:         totalTax,
        EffectiveRate:    effectiveRate,
        BracketBreakdown: breakdown,
    }, nil
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}