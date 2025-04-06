package traefik_deep_linking_middleware

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Config holds the deep linking configuration
type Config struct {
	Redirects map[string]string `json:"redirects,omitempty"`
}

// CreateConfig initializes a default configuration with common redirect mappings
func CreateConfig() *Config {
	return &Config{
		Redirects: map[string]string{
			"/validate-mail":   "app://validate-mail",
			"/change-initiate": "app://change-initiate",
			"/forgot-password": "app://forgot-password",
			"/change-email":    "app://change-email",
		},
	}
}

// DeepLinking handles deep linking redirects for mobile devices
type DeepLinking struct {
	next         http.Handler
	redirects    map[string]string
	mobileRegexp *regexp.Regexp
}

// New creates a new DeepLinking middleware instance
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if next == nil {
		return nil, ErrNilNextHandler
	}
	if config == nil {
		return nil, ErrNilConfig
	}

	mobileRE := regexp.MustCompile(`(?i)(android|blackberry|iphone|ipad|ipod|iemobile|opera mobile|webos)`)

	return &DeepLinking{
		next:         next,
		redirects:    normalizeRedirects(config.Redirects),
		mobileRegexp: mobileRE,
	}, nil
}

// ErrNilNextHandler is returned when the next handler is nil
var ErrNilNextHandler = errors.New("next handler cannot be nil")

// ErrNilConfig is returned when the config is nil
var ErrNilConfig = errors.New("config cannot be nil")

// normalizeRedirects ensures all paths start with "/" and removes trailing slashes
func normalizeRedirects(redirects map[string]string) map[string]string {
	normalized := make(map[string]string, len(redirects))
	for k, v := range redirects {
		key := strings.TrimSpace(k)
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}
		key = strings.TrimSuffix(key, "/")
		normalized[key] = strings.TrimSpace(v)
	}
	return normalized
}

// ServeHTTP implements the http.Handler interface
func (d *DeepLinking) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := strings.TrimSuffix(req.URL.Path, "/")
	ua := req.Header.Get("User-Agent")

	if ua != "" && d.mobileRegexp.MatchString(ua) {
		if redirectURL, exists := d.redirects[path]; exists && redirectURL != "" {
			http.Redirect(rw, req, redirectURL, http.StatusFound)
			return
		}
	}

	d.next.ServeHTTP(rw, req)
}
