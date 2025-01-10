package database

import "database/sql"

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Next() bool
	Close() error
	Err() error
	Scan(dest ...interface{}) error
}

type DB interface {
	QueryRow(query string, args ...interface{}) Row
	Query(query string, args ...interface{}) (Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SqlDBAdapter struct {
	*sql.DB
}

func (db *SqlDBAdapter) QueryRow(query string, args ...interface{}) Row {
	return db.DB.QueryRow(query, args...)
}

func (db *SqlDBAdapter) Query(query string, args ...interface{}) (Rows, error) {
	return db.DB.Query(query, args...)
}

func (db *SqlDBAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}
