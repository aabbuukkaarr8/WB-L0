package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New() *Store {
	return &Store{}
}

func (s *Store) Open(config string) error {
	db, err := sql.Open("postgres", config)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Store) GetConn() *sql.DB {
	return s.db
}

func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		return
	}
}
func (s *Store) SetConn(db *sql.DB) {
	s.db = db
}

func (s *Store) ConfigurePool(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) {
	if s.db == nil {
		return
	}
	if maxOpen > 0 {
		s.db.SetMaxOpenConns(maxOpen)
	}
	if maxIdle >= 0 {
		s.db.SetMaxIdleConns(maxIdle)
	}
	if maxLifetime > 0 {
		s.db.SetConnMaxLifetime(maxLifetime)
	}
	if maxIdleTime > 0 {
		s.db.SetConnMaxIdleTime(maxIdleTime)
	}
}
