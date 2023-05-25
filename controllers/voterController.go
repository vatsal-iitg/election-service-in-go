package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "github.com/vatsal-iitg/election-service-in-go/database"
	helper "github.com/vatsal-iitg/election-service-in-go/helpers"
	models "github.com/vatsal-iitg/election-service-in-go/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterVoter(c *gin.Context) {
	var voter models.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// if the constituency exists
	var constituencyID int
	err = db.QueryRow("SELECT id FROM constituencies WHERE id = $1", voter.ConstituencyID).Scan(&constituencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constituency ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check constituency"})
		}
		return
	}

	//if the candidate is valid
	if voter.CandidateID != 0 {
		var candidateID int
		err = db.QueryRow("SELECT id FROM candidate_constituencies WHERE candidate_id = $1 AND constituency_id = $2",
			voter.CandidateID, voter.ConstituencyID).Scan(&candidateID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID for the specified constituency"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check candidate"})
			}
			return
		}
	}

	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(voter.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// database insertion
	_, err = db.Exec("INSERT INTO voters (name, age, email, password, constituency_id, candidate_id) VALUES ($1, $2, $3, $4, $5, $6)",
		voter.Name, voter.Age, voter.Email, hashedPassword, voter.ConstituencyID, voter.CandidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register voter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voter registered successfully"})
}

// voter login
func LoginVoter(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON while logging in voter"})
		return
	}

	// voter info based on email
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	var voter models.Voter
	err = db.QueryRow("SELECT id, name, age, email, password FROM voters WHERE email = $1", credentials.Email).
		Scan(&voter.ID, &voter.Name, &voter.Age, &voter.Email, &voter.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	// password comparision
	err = bcrypt.CompareHashAndPassword([]byte(voter.Password), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	//  token generation
	token, err := helper.GenerateTokenForVoter(strconv.Itoa(voter.ID), voter.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
