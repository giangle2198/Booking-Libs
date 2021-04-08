package db

import "database/sql"

// DBHelper is helper of DB 
type DBHelper interface {
	Open() *sql.DB
	Close() error
	Begin() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	RollBack(tx *sql.Tx) error
}