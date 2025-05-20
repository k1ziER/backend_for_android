package service

import (
	"android/internal/domain"
	"android/internal/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.User
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) SignIn(login, password string) (domain.User, error) {
	hashPassword := s.generatePasswordHash(password)
	user, err := s.repo.SignIn(login, hashPassword)
	if err != nil {
		logrus.Println(err)
		return user, err
	}
	return user, nil
}

func (s *AuthService) GetUser(id int) (domain.User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		logrus.Println(err)
	}
	return user, err
}

func (s *AuthService) GenerateToken(user domain.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	const tokenTTL = 12 * time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(os.Getenv("SIGN_KEY")))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SIGN_KEY")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SOLT"))))
}
