package auth

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtService struct{}

func NewJwtService() Service {
	return &jwtService{}
}

var (
	secretKey  = os.Getenv("SECRET_KEY")
	SECRET_KEY = []byte(secretKey)
)

func (s *jwtService) GenerateToken(userID string) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return SECRET_KEY, nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
