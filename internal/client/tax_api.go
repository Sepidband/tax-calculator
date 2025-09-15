package client

import (
	"fmt"
	"sync"
	"tax-calculator/pkg/models"
	"time"

	"github.com/go-resty/resty/v2"
)

type TaxAPIClient struct {
	baseURL string
	client  *resty.Client
	cache   map[int]*cacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
}

type cacheEntry struct {
	data      *models.TaxBracketsResponse
	timestamp time.Time
}

func NewTaxAPIClient(baseURL string) *TaxAPIClient {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)

	return &TaxAPIClient{
		baseURL: baseURL,
		client:  client,
		cache:   make(map[int]*cacheEntry),
		ttl:     24 * time.Hour, // Cache for 24 hours
	}
}

func (c *TaxAPIClient) GetTaxBrackets(year int) (*models.TaxBracketsResponse, error) {
	// Check cache first
	if cached := c.getFromCache(year); cached != nil {
		return cached, nil
	}

	// Cache miss - fetch from API
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

	// Store in cache
	c.setCache(year, &response)

	return &response, nil
}

func (c *TaxAPIClient) getFromCache(year int) *models.TaxBracketsResponse {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[year]
	if !exists {
		return nil
	}

	// Check if cache entry is still valid
	if time.Since(entry.timestamp) > c.ttl {
		// Cache expired, remove it
		delete(c.cache, year)
		return nil
	}

	return entry.data
}

func (c *TaxAPIClient) setCache(year int, data *models.TaxBracketsResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[year] = &cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

// ClearCache removes all cached entries
func (c *TaxAPIClient) ClearCache() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[int]*cacheEntry)
}
