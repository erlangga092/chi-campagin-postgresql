package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID string) (string, error)
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
