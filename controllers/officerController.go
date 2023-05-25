package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	database "github.com/vatsal-iitg/election-service-in-go/database"
	helper "github.com/vatsal-iitg/election-service-in-go/helpers"
	models "github.com/vatsal-iitg/election-service-in-go/models"
	"golang.org/x/crypto/bcrypt"
)

// register a election officer.
func RegisterElectionOfficer(c *gin.Context) {
	log.Println("Entering RegisterElectionOfficer")

	var officer models.ElectionOfficer
	if err := c.ShouldBindJSON(&officer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("binding json completed")

	// setting up database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v in registerOfficer", err)
	}
	defer db.Close()
	log.Println("database connection completed")

	// checking the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	} else {
		log.Println("Database connection established")
	}

	// if the officer already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM election_officers WHERE email = $1", officer.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred/ can't check if the officer already exists"})
		return
	}
	log.Println("got the count")

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Election officer with the given email already exists"})
		return
	}
	log.Println("checked the count condition")

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(officer.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred/ can't hash password while registering officer"})
		return
	}
	log.Println("hashed the password")

	// officer insertion
	_, err = db.Exec("INSERT INTO election_officers (name, email, password, role) VALUES ($1, $2, $3, $4)",
		officer.Name, officer.Email, string(hashedPassword), officer.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred/ can't insert officer into the database"})
		return
	}
	log.Println("inserted into database")

	c.JSON(http.StatusOK, gin.H{"success": "officer registered successfully"})
	log.Println("Exiting RegisterElectionOfficer")
}

// officer login
func LoginElectionOfficer(c *gin.Context) {
	var credentials models.LoginCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
		return
	}

	var officer models.ElectionOfficer

	// setting up database connection
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v failed to login as officer", err)
	}
	defer db.Close()

	// getting officer info based on email
	err = db.QueryRow("SELECT id, name, email, password, role FROM election_officers WHERE email = $1 AND role = $2",
		credentials.Email, credentials.Role).
		Scan(&officer.ID, &officer.Name, &officer.Email, &officer.Password, &officer.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials for login officer/ no such officer"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred/ error while logging in for officer"})
		}
		return
	}

	// password verification
	err = bcrypt.CompareHashAndPassword([]byte(officer.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials/ password dont match for officer login"})
		return
	}

	// token generation
	token, err := helper.GenerateTokenForOfficer(officer.ID, officer.Email, officer.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred/ can't get token for officer login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
