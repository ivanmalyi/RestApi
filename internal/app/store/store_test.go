package store_test

import (
	"os"
	"testing"
)

var databaseUrl string

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		databaseUrl = "postgresql://root:1@localhost:5432/venera_test"
	}

	os.Exit(m.Run())
}