package main

import (
	"coding-challenge/routers"
	"log"
	"net/http"
	"os"
)

func main() {
	/* Set Env it will be accessable anywhere in app
	a) For production use keyword "Production"
	b) For staging use keyword "Staging"
	c) For development use keyword "Development" */

	os.Setenv("ENV", "Development")
	//Env default set
	env := os.Getenv("ENV")

	//Check env and set the port
	var port string
	if env == "Development" {
		port = "4000"
		// Log server started
		log.Println("Server started at port ", port)
	}

	// Initalize all the routes and start the server
	router := routers.InitRoutes()
	httpError := http.ListenAndServe(":"+port, router)
	if httpError != nil {
		log.Println("While serving HTTP: ", httpError)
	}
}
