package application

import (
	"encoding/json"
	"sync"
)

// BrowserDataStore is a thread-safe store for browser data extracted from browser mode windows.
// It allows sharing extracted data between windows and provides methods to access the data
// from Go code.
type BrowserDataStore struct {
	mu   sync.RWMutex
	data map[string]*BrowserData // keyed by window name
}

// NewBrowserDataStore creates a new BrowserDataStore instance
func NewBrowserDataStore() *BrowserDataStore {
	return &BrowserDataStore{
		data: make(map[string]*BrowserData),
	}
}

// globalBrowserDataStore is the global instance of the browser data store
var globalBrowserDataStore = NewBrowserDataStore()

// GetBrowserDataStore returns the global browser data store instance
func GetBrowserDataStore() *BrowserDataStore {
	return globalBrowserDataStore
}

// Store stores browser data for a window
func (s *BrowserDataStore) Store(windowName string, data *BrowserData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[windowName] = data
}

// Get retrieves browser data for a window
func (s *BrowserDataStore) Get(windowName string) *BrowserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[windowName]
}

// GetAll returns all stored browser data
func (s *BrowserDataStore) GetAll() map[string]*BrowserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]*BrowserData)
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

// Delete removes browser data for a window
func (s *BrowserDataStore) Delete(windowName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, windowName)
}

// Clear removes all stored browser data
func (s *BrowserDataStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]*BrowserData)
}

// GetCookies returns cookies for a specific window
func (s *BrowserDataStore) GetCookies(windowName string) map[string]string {
	data := s.Get(windowName)
	if data == nil {
		return nil
	}
	return data.Cookies
}

// GetLocalStorage returns localStorage for a specific window
func (s *BrowserDataStore) GetLocalStorage(windowName string) map[string]string {
	data := s.Get(windowName)
	if data == nil {
		return nil
	}
	return data.LocalStorage
}

// GetSessionStorage returns sessionStorage for a specific window
func (s *BrowserDataStore) GetSessionStorage(windowName string) map[string]string {
	data := s.Get(windowName)
	if data == nil {
		return nil
	}
	return data.SessionStorage
}

// GetURLParams returns URL parameters for a specific window
func (s *BrowserDataStore) GetURLParams(windowName string) map[string]string {
	data := s.Get(windowName)
	if data == nil {
		return nil
	}
	return data.URLParams
}

// GetHTMLContent returns HTML content for a specific window
func (s *BrowserDataStore) GetHTMLContent(windowName string) string {
	data := s.Get(windowName)
	if data == nil {
		return ""
	}
	return data.HTMLContent
}

// GetCustomData returns custom data for a specific window
func (s *BrowserDataStore) GetCustomData(windowName string) map[string]interface{} {
	data := s.Get(windowName)
	if data == nil {
		return nil
	}
	return data.CustomData
}

// ToJSON returns the browser data as a JSON string
func (d *BrowserData) ToJSON() (string, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseBrowserDataFromJSON parses browser data from a JSON string
func ParseBrowserDataFromJSON(jsonStr string) (*BrowserData, error) {
	var data BrowserData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

