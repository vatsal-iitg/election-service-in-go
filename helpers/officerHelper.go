package helper

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenForOfficer(id int, email string, role string) (string, error) {

	//expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   expirationTime.Unix(),
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load the environment file: %v", err)
		return "", err
	}

	// craeting the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing the token with the jwt secret key

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
