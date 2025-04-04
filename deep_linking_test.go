package deep_linking

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		next    http.Handler
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			next:    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			config:  CreateConfig(),
			wantErr: false,
		},
		{
			name:    "nil next handler",
			next:    nil,
			config:  CreateConfig(),
			wantErr: true,
		},
		{
			name:    "nil config",
			next:    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(context.Background(), tt.next, tt.config, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMobileRedirect_ServeHTTP(t *testing.T) {
	config := CreateConfig()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("next handler"))
	})

	middleware, err := New(context.Background(), nextHandler, config, "test")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		path         string
		userAgent    string
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{
			name:         "mobile with redirect",
			path:         "/validate-mail",
			userAgent:    "iPhone",
			wantStatus:   http.StatusFound,
			wantLocation: "app://validate-mail",
		},
		{
			name:       "mobile with unknown path",
			path:       "/unknown",
			userAgent:  "Android",
			wantStatus: http.StatusOK,
			wantBody:   "next handler",
		},
		{
			name:       "desktop with redirect path",
			path:       "/forgot-password",
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantStatus: http.StatusOK,
			wantBody:   "next handler",
		},
		{
			name:         "mobile with trailing slash",
			path:         "/change-email/",
			userAgent:    "iPad",
			wantStatus:   http.StatusFound,
			wantLocation: "app://change-email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("User-Agent", tt.userAgent)
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("ServeHTTP() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantLocation != "" {
				location := rr.Header().Get("Location")
				if location != tt.wantLocation {
					t.Errorf("ServeHTTP() location = %v, want %v", location, tt.wantLocation)
				}
			}

			if tt.wantBody != "" && rr.Body.String() != tt.wantBody {
				t.Errorf("ServeHTTP() body = %v, want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestNormalizeRedirects(t *testing.T) {
	input := map[string]string{
		"validate-email": "app://validate-email",
		"/change-email/": "app://change-email",
		"  /space  ":     "app://space",
	}

	expected := map[string]string{
		"/validate-email": "app://validate-email",
		"/change-email":   "app://change-email",
		"/space":          "app://space",
	}

	result := normalizeRedirects(input)

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("normalizeRedirects() = %v, want %v for key %v", result[k], v, k)
		}
	}
}
