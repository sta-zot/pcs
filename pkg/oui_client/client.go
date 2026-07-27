package ouiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Provider определяет интерфейс для любого сервиса поиска OUI.
// Это позволяет нам использовать единый цикл для разных API (паттерн Адаптер).
type Provider interface {
	Name() string
	Fetch(ctx context.Context, mac string) (string, error)
}

// --- Структуры для парсинга JSON ---

// MacVendorLookupResponse описывает ответ от macvendorlookup.com (массив)
type MacVendorLookupResponse []struct {
	Company string `json:"company"`
}

// MacLookupAppResponse описывает ответ от maclookup.app (объект)
type MacLookupAppResponse struct {
	Success bool   `json:"success"`
	Found   bool   `json:"found"`
	Company string `json:"company"`
}

// --- Реализации провайдеров ---

type MacVendorLookupProvider struct {
	client *http.Client
}

func (p *MacVendorLookupProvider) Name() string { return "macvendorlookup" }

func (p *MacVendorLookupProvider) Fetch(ctx context.Context, mac string) (string, error) {
	url := fmt.Sprintf("https://www.macvendorlookup.com/api/v2/%s", mac)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data MacVendorLookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data) > 0 && data[0].Company != "" {
		return data[0].Company, nil
	}
	return "", fmt.Errorf("no data found")
}

type MacLookupAppProvider struct {
	client *http.Client
	// apiKey string // Можно добавить для платных лимитов
}

func (p *MacLookupAppProvider) Name() string { return "maclookup" }

func (p *MacLookupAppProvider) Fetch(ctx context.Context, mac string) (string, error) {
	url := fmt.Sprintf("https://api.maclookup.app/v2/macs/%s", mac)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data MacLookupAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Success && data.Found && data.Company != "" {
		return data.Company, nil
	}
	return "", fmt.Errorf("no data found")
}

// --- Основной клиент ---

type OUIClient struct {
	providers  []Provider
	httpClient *http.Client

	// Простой in-memory кэш для снижения нагрузки на API
	cache map[string]string
	mu    sync.RWMutex

	// Регулярное выражение для очистки юридических суффиксов в конце строки
	// Удаляет: INC, LTD, LLC, CORP, Co., Ltd, GMBH и т.д.
	cleanupRegex *regexp.Regexp
}

func NewOUIClient() *OUIClient {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &OUIClient{
		providers: []Provider{
			&MacVendorLookupProvider{client: httpClient},
			&MacLookupAppProvider{client: httpClient},
		},
		httpClient:   httpClient,
		cache:        make(map[string]string),
		cleanupRegex: regexp.MustCompile(`(?i)\s*,?\s*\b(?:inc\.?|ltd\.?|llc\.?|corp\.?|corporation|company|co\.?,?\s*ltd\.?|gmbh|ag|sa|bv|ab|oy|pte|pty|limited|computer|technology|electronics|Enterprise|(?:\([a-zA-Z0-9]*\)))\.?$`),
	}
}

// normalizeMac очищает MAC-адрес от разделителей и приводит к нижниму регистру
func normalizeMac(mac string) string {
	// Оставляем только hex-символы
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	clean := re.ReplaceAllString(mac, "")
	return strings.ToLower(clean)
}

// cleanCompanyName удаляет юридические суффиксы из названия
func (c *OUIClient) cleanCompanyName(name string) string {
	clean := ""
	for i := 0; i < len(name); i++ {
		clean = strings.TrimSpace(c.cleanupRegex.ReplaceAllString(name, ""))
		if clean == strings.TrimSpace(c.cleanupRegex.ReplaceAllString(clean, "")) {
			break
		}
		name = clean
	}
	return clean
}

// GetVendor ищет производителя, перебирая провайдеров до первого успешного результата
func (c *OUIClient) GetVendor(ctx context.Context, mac string) (string, error) {
	cleanMac := normalizeMac(mac)

	// 1. Проверка кэша
	c.mu.RLock()
	if cached, ok := c.cache[cleanMac]; ok {
		c.mu.RUnlock()
		if cached == "" {
			return "", fmt.Errorf("manufacturer not found (cached)")
		}
		return cached, nil
	}
	c.mu.RUnlock()

	// 2. Перебор провайдеров (Fallback механизм)
	var lastErr error
	for _, provider := range c.providers {
		company, err := provider.Fetch(ctx, cleanMac)
		if err == nil && company != "" {
			cleanedCompany := c.cleanCompanyName(company)

			// Сохраняем в кэш
			c.mu.Lock()
			c.cache[cleanMac] = cleanedCompany
			c.mu.Unlock()

			return cleanedCompany, nil
		}
		lastErr = err
		// Логируем ошибку и переходим к следующему провайдеру
		// fmt.Printf("[WARN] Provider %s failed: %v\n", provider.Name(), err)
	}

	// Если ни один провайдер не справился, кэшируем пустой результат, чтобы не долбить API
	// c.mu.Lock()
	// c.cache[cleanMac] = ""
	// c.mu.Unlock()

	return "", fmt.Errorf("all providers failed. last error: %w", lastErr)
}
