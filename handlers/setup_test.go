package handlers

import (
	"os"
	"testing"

	"github.com/carlosarraes/fsback/repository/dbrepo"
)

var app App

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	os.Exit(m.Run())
}
