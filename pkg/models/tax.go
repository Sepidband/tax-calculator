package models

type TaxBracket struct {
    Min  float64 `json:"min"`
    Max  *float64 `json:"max,omitempty"` // pointer for unlimited top bracket
    Rate float64 `json:"rate"`
}

type TaxBracketsResponse struct {
    TaxBrackets []TaxBracket `json:"tax_brackets"`
}

type TaxCalculationResult struct {
    TotalTax        float64            `json:"total_tax"`
    EffectiveRate   float64            `json:"effective_rate"`
    BracketBreakdown []BracketTaxDetail `json:"bracket_breakdown"`
}

type BracketTaxDetail struct {
    Bracket TaxBracket `json:"bracket"`
    TaxOwed float64    `json:"tax_owed"`
}