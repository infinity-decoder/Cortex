package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtKey []byte

func init() {
	jwtKeyStr := os.Getenv("JWT_SECRET_KEY")
	if jwtKeyStr == "" {
		// Generate random 32-byte key
		jwtKey = make([]byte, 32)
		if _, err := rand.Read(jwtKey); err != nil {
			log.Fatalf("Failed to generate JWT secret key: %v", err)
		}
		log.Println("WARNING: JWT_SECRET_KEY not set, using random key (tokens will be invalid on restart)")
	} else {
		if len(jwtKeyStr) < 32 {
			log.Fatalf("JWT_SECRET_KEY must be at least 32 characters long")
		}
		jwtKey = []byte(jwtKeyStr)
	}
}

// ValidateJWTSecret validates that JWT secret is properly configured
func ValidateJWTSecret() error {
	if len(jwtKey) < 32 {
		return fmt.Errorf("JWT secret key is too short (minimum 32 bytes)")
	}
	return nil
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	OrgID  uuid.UUID `json:"org_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, orgID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		UserID: userID,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
