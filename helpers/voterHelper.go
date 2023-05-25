package helper

import (
	"os"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

const (
	voterTokenExpirationTime = 24 * time.Hour // token expiration time
)

// jwt token for voter
func GenerateTokenForVoter(voterID, email string) (string, error) {

	// token claims
	claims := jwt.MapClaims{
		"voterID": voterID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
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
