package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	initSQLDatabase(cfg.Database)
}

func initSQLDatabase(dbCfg config.DatabaseConfig) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
		dbCfg.SslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}
