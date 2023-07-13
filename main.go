package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	// for loading data from env file
	"github.com/joho/godotenv"

	// importing package routes
	routes "github.com/vatsal-iitg/election-service-in-go/routes"
)

//  gin contains a set of commonly used functionalities (e.g., routing, middleware support, rendering, etc.) that reduce boilerplate code and make it simpler to build web applications

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load the environment file: %v", err)
		return
	}

	port := ""
	port = os.Getenv("PORT")
	if port == "" {
		port = "5400"
	}

	router := gin.New() // creating a new instance of the gin router
	// In Go, a router is responsible for routing incoming HTTP requests to the appropriate handler functions based on the requested URL and HTTP method. Gin provides a convenient router implementation that simplifies this process.

	// The gin.New() function creates a new instance of the Gin router. This function initializes the router with default settings and returns a pointer to the router object. By assigning the returned router to the variable router, you can then configure routes and define handlers for different endpoints using the router.
	routes.RouterHandler(router)

	router.Run(":" + port)
	// used to start the HTTP server and listen for incoming requests on a specific port.

}
