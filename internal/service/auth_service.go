package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"smartqueue/internal/models"
	"smartqueue/internal/repository"
)

type AuthService struct {
	users     *repository.UserPostgresRepository
	jwtSecret []byte
}

type AuthClaims struct {
	UserID int    `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

func NewAuthService(users *repository.UserPostgresRepository, jwtSecret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
	}
}

func normalizePhone(phone string) string {
	var b strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

var validPriorityCategories = map[string]bool{
	models.PriorityCategoryNone:     true,
	models.PriorityCategoryPregnant: true,
	models.PriorityCategoryElderly:  true,
	models.PriorityCategoryDisabled: true,
}

func (s *AuthService) Register(firstName, lastName, phone, password, priorityCategory string) (*models.User, string, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	phone = normalizePhone(phone)
	password = strings.TrimSpace(password)

	if priorityCategory == "" {
		priorityCategory = models.PriorityCategoryNone
	}

	if firstName == "" || lastName == "" || phone == "" || password == "" {
		return nil, "", errors.New("all fields are required")
	}
	if len(password) < 6 {
		return nil, "", errors.New("password must be at least 6 characters")
	}
	if !validPriorityCategories[priorityCategory] {
		return nil, "", fmt.Errorf("invalid priority_category: must be one of none, pregnant, elderly, disabled")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		FirstName:        firstName,
		LastName:         lastName,
		Phone:            phone,
		PasswordHash:     string(hash),
		PriorityCategory: priorityCategory,
	}

	if err := s.users.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, "", err
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *AuthService) Login(phone, password string) (*models.User, string, error) {
	phone = normalizePhone(phone)
	password = strings.TrimSpace(password)

	if phone == "" || password == "" {
		return nil, "", errors.New("phone and password are required")
	}

	user, err := s.users.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("invalid phone or password")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid phone or password")
	}

	token, err := s.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, "", err
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	user, err := s.users.GetByID(id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) GenerateToken(userID int, phone string) (string, error) {
	now := time.Now()

	claims := AuthClaims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "smartqueue",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
