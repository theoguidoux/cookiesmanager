// Package cookiesmanager traefik middleware plugin.
package cookiesmanager

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type CookieConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Config the plugin configuration.
type Config struct {
	Adder   []CookieConfig `json:"adder,omitempty"`
	Remover []CookieConfig `json:"remover,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Adder:   []CookieConfig{},
		Remover: []CookieConfig{},
	}
}

type CookieManager struct {
	next    http.Handler
	adder   []CookieConfig
	remover []CookieConfig
	name    string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CookieManager{
		adder:   config.Adder,
		remover: config.Remover,
		next:    next,
		name:    name,
	}, nil
}

func (c *CookieManager) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Copy existing cookies
	cookies := req.Cookies()

	// Prepare maps of cookies to easily access them
	removerMap := make(map[string]string)
	for _, v := range c.remover {
		removerMap[v.Name] = v.Value
	}

	adderMap := make(map[string]string)
	for _, v := range c.adder {
		adderMap[v.Name] = v.Value
	}

	// Remove all cookies from request
	req.Header.Set("Cookie", "")

	// Add and Remove cookies content
	for _, cookie := range cookies {

		removingCookieValue, rmOk := removerMap[cookie.Name]
		addingCookieValue, addOk := adderMap[cookie.Name]

		if addOk {
			if !strings.Contains(cookie.Value, addingCookieValue) {
				cookie.Value = fmt.Sprintf("%s %s", cookie.Value, addingCookieValue)
			}
		}

		if rmOk {
			if strings.Contains(cookie.Value, removingCookieValue) {
				cookie.Value = strings.ReplaceAll(cookie.Value, removingCookieValue, "")
			}
		}

		req.AddCookie(cookie)
	}

	// Create missing cookie if does not exist the request
	for _, cookie := range c.adder {

		_, err := req.Cookie(cookie.Name)

		if err == http.ErrNoCookie {
			foundCookie := &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			}

			req.AddCookie(foundCookie)
		}
	}

	c.next.ServeHTTP(rw, req)
}
