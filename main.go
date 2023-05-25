package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	routes "github.com/vatsal-iitg/election-service-in-go/routes"
)

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

	router := gin.New()
	routes.RouterHandler(router)

	router.Run(":" + port)

}
