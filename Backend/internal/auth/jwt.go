package auth
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateToken(userId string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		return "", errors.New("JWT_SECRET environment variable is not configured")
	}
	jwtSecret := []byte(secretKey)
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	// 1. Fetch the secret key from environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable is not configured")
	}
	jwtSecret := []byte(secretKey)

	// 2. Parse the token with our custom Claims struct
	// The inline callback function passes your secret key back to the internal parser.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		// ⚠️ CRITICAL SECURITY STEP: Validate the signing algorithm
		// This guarantees the token was signed with HS256.
		// It prevents a known exploit where hackers pass 'alg: "none"' to bypass verification.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	// 3. Handle parsing errors (e.g., malformed tokens, bad characters)
	if err != nil {
		return nil, err
	}

	// 4. Extract and check validity
	// token.Valid automatically checks if the token has expired (exp claim)
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token or expired signatures")
}

func main() {
	// userId := "wfhohh2irj"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoid2Zob2hoMmlyaiIsImV4cCI6MTc4MTg1OTk1NiwiaWF0IjoxNzgxNzczNTU2fQ.ZWOCEYMDC_UdxK9q0sJVrLs29N6G9k8m-JV7mKaFp70"
	claims, err := VerifyToken(token)
	if err != nil {
		log.Fatalf("Failed to verify token: %v", err)
	}

	payload, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format payload: %v", err)
	}

	fmt.Printf("Verified JWT payload:\n%s\n", payload)
}


