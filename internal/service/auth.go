package service

import (
	"android/pkg/domain"
	"android/pkg/ports"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type tokenTicketClaims struct {
	jwt.StandardClaims
	UserId          int       `json:"user_id"`
	TitleAttraction string    `json:"title_attraction"`
	Date            time.Time `json:"date"`
	Count           int       `json:"count"`
}

type UserBlackList struct {
	UserId map[int]bool
	mu     sync.RWMutex
	Token  map[string]bool
}
type AuthService struct {
	repo      ports.UserRepo
	blackList ports.UserBlackList
}

func NewTokenBlacklist() *UserBlackList {
	return &UserBlackList{
		UserId: make(map[int]bool),
		Token:  make(map[string]bool),
	}
}

func NewAuthService(repo ports.UserRepo, blackList ports.UserBlackList) *AuthService {
	return &AuthService{
		repo:      repo,
		blackList: blackList,
	}
}

func (s *AuthService) CreateUser(user domain.User) (domain.User, error) {
	password, err := s.validatePassword(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = s.generatePasswordHash(password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) UpdateUser(user domain.User) error {
	return s.repo.UpdateUser(user)
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

func (s *AuthService) DeleteUser(id int) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		logrus.Println(err)
	}
	s.blackList.AddUserBlackList(id)

	return err
}

func (s *AuthService) Logout(token string) {
	s.blackList.AddTokenBlackList(token)
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

func (s *AuthService) CreateTicket(id int, ticket domain.Ticket) (string, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	const tokenTTL = 3 * time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenTicketClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:          id,
		TitleAttraction: ticket.TitleAttraction,
		Date:            ticket.Date,
		Count:           ticket.Count,
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

func (s *AuthService) ParseTicketToken(accessToken string) (domain.Ticket, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}
	token, err := jwt.ParseWithClaims(accessToken, &tokenTicketClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SIGN_KEY")), nil
	})
	if err != nil {
		return domain.Ticket{}, err
	}

	claims, ok := token.Claims.(*tokenTicketClaims)
	if !ok {
		return domain.Ticket{}, errors.New("token claims are not of type *tokenClaims")
	}

	return domain.Ticket{
		TitleAttraction: claims.TitleAttraction,
		Date:            claims.Date,
		Count:           claims.Count,
	}, nil
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

func (s *AuthService) validatePassword(password string) (string, error) {
	hasLower := false
	hasUpper := false
	hasSpace := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsSpace(char):
			hasSpace = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if hasLower == true && hasDigit == true && hasSpace == false && hasUpper == true && hasSpecial == true {
		return password, nil
	}

	return "", errors.New("password is unreliable")
}

func (tb *UserBlackList) AddUserBlackList(userId int) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.UserId[userId] = true
}

func (tb *UserBlackList) AddTokenBlackList(token string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.Token[token] = true
}

func (tb *UserBlackList) IsUserBlackListed(userId int) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.UserId[userId]
}

func (tb *UserBlackList) IsTokenBlackListed(token string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.Token[token]
}
