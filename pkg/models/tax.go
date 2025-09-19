// pkg/models/tax.go
package models

// TaxBracket represents a tax bracket with minimum income, maximum income, and tax rate
type TaxBracket struct {
	Min  float64  `json:"min"`
	Max  *float64 `json:"max,omitempty"` // pointer for unlimited top bracket
	Rate float64  `json:"rate"`
}

// TaxBracketsResponse represents the response from the external tax API
type TaxBracketsResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

// TaxCalculationResult represents the complete tax calculation result
type TaxCalculationResult struct {
	TotalTax         float64            `json:"total_tax"`
	EffectiveRate    float64            `json:"effective_rate"`
	BracketBreakdown []BracketTaxDetail `json:"bracket_breakdown"`
}

// BracketTaxDetail represents tax calculation details for a specific bracket
type BracketTaxDetail struct {
	Bracket TaxBracket `json:"bracket"`
	TaxOwed float64    `json:"tax_owed"`
}
