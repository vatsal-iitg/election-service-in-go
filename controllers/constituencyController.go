package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/vatsal-iitg/election-service-in-go/database"
	models "github.com/vatsal-iitg/election-service-in-go/models"
)

// creating the constituency
func CreateConstituency(c *gin.Context) {
	var constituency models.Constituency
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json not binded"})
		return
	}
	log.Println("binded json for constituency creation")

	// database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v in createConstituency", err)
		return
	}
	defer db.Close()
	log.Println("connection established with database")

	// database insertion
	err = db.QueryRow("INSERT INTO constituencies (name) VALUES ($1) RETURNING id", constituency.Name).Scan(&constituency.ID)
	if err != nil {
		log.Fatalf("Failed to insert constituency: %v in createConstituency", err)
		return
	}

	c.JSON(http.StatusCreated, constituency)
	log.Println("constituency created")
}

// updating the constituency
func UpdateConstituency(c *gin.Context) {

	var constituency models.UpdateConstituencyCredentials
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v in createConstituency", err)
		return
	}
	defer db.Close()

	// updating the data
	_, err = db.Exec("UPDATE constituencies SET name = $1 WHERE id = $2", constituency.Name, constituency.ID)
	if err != nil {
		log.Fatalf("Failed to update constituency: %v", err)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update constituency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Constituency updated successfully"})
}
