package database

import (
	"time"

	CONFIG "merge-backend/config"

	"github.com/google/uuid"
)

var MainDB Postgres

// ConnectDatabases - connects all databases with given configurations
func ConnectDatabases() error {
	// connect main postgres db
	MainDB = Postgres{dbConfig: CONFIG.MainDBConfig, connectionPool: CONFIG.MainDBConnectionPool, maxLifeTime: time.Hour}
	err := MainDB.Connect()
	if err != nil {
		return err
	}

	return nil
}

func generateRandomID() string {
	return uuid.New().String()
}
