package tnserver

import (
	"database/sql"
	"net/http"

	"github.com/artemxgod/transaction-server/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
)

func Start(cfg *Config) error {
	db, err := newDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)

	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionKey))
	s := newServer(store, sessionStore)

	return http.ListenAndServe(cfg.BindAddr, s)
}

func newDB(database_url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", database_url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}