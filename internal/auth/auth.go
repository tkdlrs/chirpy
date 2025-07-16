package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT -
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	//
	mySigningKey := []byte("AllYourBase")
	// claims
	currentTime := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTime.UTC()),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
		Subject:   userID.String(),
	}

	// token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signing
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	// happy path
	return tokenString, nil
}

// ValidateJWT
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	//
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}
	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//
	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		subject, err := claims.GetSubject()
		if err != nil {
			return uuid.UUID{}, err
		}
		return uuid.MustParse(subject), nil
	}
	//
	return uuid.UUID{}, nil
}
