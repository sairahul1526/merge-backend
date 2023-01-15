package main

import (
	API "merge-backend/api"
	INIT "merge-backend/init"
)

func main() {

	INIT.Init()

	// start api server
	API.StartServer()

}
