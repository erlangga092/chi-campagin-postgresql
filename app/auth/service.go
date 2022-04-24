package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"funding-app/app/key"
	"funding-app/app/user"
	"io"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	GenerateToken(userID string) (key.Token, error)
	GenerateRefreshToken(token key.Token) (key.Token, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
	ValidateRefreshToken(token key.Token) (user.User, error)
}

type jwtService struct{}

func NewJwtService() Service {
	return &jwtService{}
}

var (
	secretKey  = os.Getenv("SECRET_KEY")
	SECRET_KEY = []byte(secretKey)
)

func (s *jwtService) GenerateToken(userID string) (key.Token, error) {
	var err error

	claim := jwt.MapClaims{}
	claim["user_id"] = userID
	claim["exp"] = time.Now().Add(time.Second * 10).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtToken := key.Token{}

	jwtToken.AccessToken, err = token.SignedString(SECRET_KEY)
	if err != nil {
		return jwtToken, err
	}

	return s.GenerateRefreshToken(jwtToken)
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

func (s *jwtService) GenerateRefreshToken(token key.Token) (key.Token, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, os.Getenv("SECRET_KEY"))

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		log.Error(err.Error())
		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))
	return token, nil
}

func (s *jwtService) ValidateRefreshToken(token key.Token) (user.User, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, os.Getenv("SECRET_KEY"))

	user := user.User{}
	// jwtToken := key.Token{}

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return user, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return user, err
	}

	data, err := base64.URLEncoding.DecodeString(token.RefreshToken)
	if err != nil {
		return user, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return user, err
	}

	if string(plain) != token.AccessToken {
		return user, errors.New("invalid token")
	}

	validatedToken, err := s.ValidateToken(string(plain))
	if err != nil {
		return user, err
	}

	claim, ok := validatedToken.Claims.(jwt.MapClaims)
	if !ok || !validatedToken.Valid {
		return user, err
	}

	user.ID = fmt.Sprintf("%s", claim["user_id"])

	return user, nil
}
