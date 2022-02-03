package appserver

import (
	"database/sql"
	"github.com/ivanmalyi/RestApi/internal/app/store/sqlstore"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	server := NewServer(store)
	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}