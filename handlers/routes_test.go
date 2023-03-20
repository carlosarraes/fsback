package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/carlosarraes/fsback/repository"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	registered := []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/users", "GET"},
		{"/user", "POST"},
		{"/user/{lastName}", "DELETE"},
	}

	mux := app.Routes()

	chiRoutes := mux.(chi.Routes)

	for _, r := range registered {
		if !checkRoute(r.route, r.method, chiRoutes) {
			t.Errorf("Route %s %s not found", r.method, r.route)
		}
	}
}

func TestConnect(t *testing.T) {
	var mock repository.DbRepo
	app := &App{DSN: "test-dsn", DB: mock}

	mockDb, err := app.Connect()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if mockDb != nil {
		t.Error("Expected nil, got a DB connection")
	}
}

func checkRoute(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false
	_ = chi.Walk(chiRoutes, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(route, testRoute) && strings.EqualFold(method, testMethod) {
			found = true
		}
		return nil
	})
	return found
}
