package client

import (
    "fmt"
    "time"
    "github.com/go-resty/resty/v2"
    "tax-calculator/pkg/models"
)

type TaxAPIClient struct {
    baseURL string
    client  *resty.Client
}

func NewTaxAPIClient(baseURL string) *TaxAPIClient {
    client := resty.New()
    client.SetTimeout(10 * time.Second)
    client.SetRetryCount(3)
    client.SetRetryWaitTime(1 * time.Second)
    
    return &TaxAPIClient{
        baseURL: baseURL,
        client:  client,
    }
}

func (c *TaxAPIClient) GetTaxBrackets(year int) (*models.TaxBracketsResponse, error) {
    var response models.TaxBracketsResponse
    
    resp, err := c.client.R().
        SetResult(&response).
        Get(fmt.Sprintf("%s/tax-calculator/tax-year/%d", c.baseURL, year))
    
    if err != nil {
        return nil, fmt.Errorf("failed to fetch tax brackets: %w", err)
    }
    
    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode())
    }
    
    return &response, nil
}