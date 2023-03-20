package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHealth(t *testing.T) {
	theTests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{"CheckHealth", "/", http.StatusOK},
		{"404", "/notfound", http.StatusNotFound},
	}

	mux := app.Routes()

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	for _, tt := range theTests {
		resp, err := ts.Client().Get(ts.URL + tt.url)
		if err != nil {
			t.Errorf("Error: %s", err)
		}

		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
		}
	}
}
