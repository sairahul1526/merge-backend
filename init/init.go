package init

import (
	"fmt"
	"math/rand"
	CONFIG "merge-backend/config"
	DB "merge-backend/database"
	"time"
)

// Init - init all config, connect databases etc
func Init() {

	rand.Seed(time.Now().UnixNano()) // seed for random generator

	// load config
	err := CONFIG.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config", err)
		return
	}

	// connect to databases
	DB.ConnectDatabases()
}
