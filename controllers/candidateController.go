package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	database "github.com/vatsal-iitg/election-service-in-go/database"
	models "github.com/vatsal-iitg/election-service-in-go/models"
)

func RegisterCandidate(c *gin.Context) {
	var candidate models.Candidate

	// bind the data as json
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// validating the data by struct
	validate := validator.New()
	if err := validate.Struct(candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// connecting to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	// is the candidate already registered
	existingCandidate, err := GetCandidateByEmail(candidate.Email)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check candidate"})
		return
	}
	if existingCandidate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Candidate already registered"})
		return
	}

	// data insertion
	insertCandidateQuery := `
		INSERT INTO candidates (name, age, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var candidateID int
	err = db.QueryRow(insertCandidateQuery, candidate.Name, candidate.Age, candidate.Email, candidate.Password).Scan(&candidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register candidate"})
		return
	}

	// mapping candidates to constituencies
	insertCandidateConstituencyQuery := `
		INSERT INTO candidate_constituencies (candidate_id, constituency_id)
		VALUES ($1, $2)
	`
	for _, constituencyID := range candidate.Constituencies {
		_, err = db.Exec(insertCandidateConstituencyQuery, candidateID, constituencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register candidate"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Candidate registered successfully"})
}

// retrieve candidate by email
func GetCandidateByEmail(email string) (*models.Candidate, error) {

	// connecting to database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// finding the candidate
	query := "SELECT id, name, email, password FROM candidates WHERE email = $1"
	row := db.QueryRow(query, email)

	// struct initialization
	candidate := &models.Candidate{}
	err = row.Scan(&candidate.ID, &candidate.Name, &candidate.Email, &candidate.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// does not exist
			return nil, nil
		}
		return nil, err
	}

	return candidate, nil
}
