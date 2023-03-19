package handlers

import (
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
