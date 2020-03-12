package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

// MemStatsLoaderInterface is the interface to load memory statistics data
type MemStatsLoaderInterface interface {
	LoadMemStats() (*runtime.MemStats, error)
}

// MemStatsLoader is a struct implements MemStatsLoaderInterface
type MemStatsLoader struct {
	url    string
	client *http.Client
}

// NewMemStatsLoader create a loader
func NewMemStatsLoader(url string) *MemStatsLoader {
	return &MemStatsLoader{
		url: url,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// LoadMemStats loads data from target application via HTTP
func (l *MemStatsLoader) LoadMemStats() (*runtime.MemStats, error) {
	resp, err := l.client.Get(l.url)
	if err != nil {
		return nil, fmt.Errorf("load memstat connect err %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("load memstat bad code err, %d", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("load memstat read  err %w", err)
	}

	// define an anonymous struct for JSON parsing
	result := &struct {
		Stats *runtime.MemStats `json:"memstats"`
	}{}
	if err := json.Unmarshal(b, result); err != nil {
		return nil, fmt.Errorf("fetch memstat, json  err %w", err)
	}
	return result.Stats, nil
}
