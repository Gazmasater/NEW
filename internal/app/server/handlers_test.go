package server

import (
	"database/sql"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/mattn/go-sqlite3"
)

func TestDBConnection(t *testing.T) {
	// Создание временной базы данных SQLite для тестирования подключения
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		t.Fatalf("error connecting to the database: %s", err)
	}
}
