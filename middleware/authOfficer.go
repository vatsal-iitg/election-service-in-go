package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	models "github.com/vatsal-iitg/election-service-in-go/models"
)

// IsAuthorizedElectionOfficer checks if the user is an authorized election officer
func IsAuthorizedElectionOfficer(c *gin.Context) bool {
	user, exists := c.Get("user")
	if !exists {
		return false
	}

	// Check if the user is an election officer
	electionOfficer, ok := user.(models.ElectionOfficer)
	if !ok {
		return false
	}

	// Perform additional authorization checks for election officer, if required
	if electionOfficer.Role != "admin" {
		return false
	}

	// Return true if all checks pass
	return true
}

// GetAuthenticatedOfficerID extracts the authenticated officer ID from the context
func GetAuthenticatedOfficerID(c *gin.Context) (string, error) {
	officerID, ok := c.Get("officer_id")
	if !ok {
		return "", errors.New("officer ID not found in context")
	}
	id, ok := officerID.(string)
	if !ok {
		return "", errors.New("invalid officer ID type")
	}
	return id, nil
}

func AuthMiddlewareForOfficer() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Failed to load the environment file: %v", err)
			return
		}

		// election officer routes
		if strings.HasPrefix(c.Request.URL.Path, "/election-officer") {
			// Get the token from the Authorization header
			tokenString := c.GetHeader("Authorization")

			// open routes
			if c.Request.URL.Path == "/election-officer/register" || c.Request.URL.Path == "/election-officer/login" {
				c.Next()
				return
			}

			// token verification for other routes
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header in authMiddleware"})
				c.Abort()
				return
			}

			// claims extraction
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
				return secretKey, nil
			})
			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token in claims"})
				c.Abort()
				return
			}
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}

func AuthorizedOnlyForOfficer() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized/ user is not admin"})
			c.Abort()
			return
		}

		c.Next()
	}
}
