package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"myapp-auth/internal/models"
	"myapp-auth/internal/repository"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Users     *repository.UserRepository
	JWTSecret []byte
}

type AuthClaims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

func NewAuthService(users *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		Users:     users,
		JWTSecret: []byte(jwtSecret),
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

func (s *AuthService) Register(ctx context.Context, firstName, lastName, phone, password string) (*models.User, string, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	phone = normalizePhone(phone)
	password = strings.TrimSpace(password)

	if firstName == "" || lastName == "" || phone == "" || password == "" {
		return nil, "", errors.New("all fields are required")
	}

	if len(password) < 6 {
		return nil, "", errors.New("password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		PasswordHash: string(hash),
	}

	if err := s.Users.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, "", err
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, phone, password string) (*models.User, string, error) {
	phone = normalizePhone(phone)
	password = strings.TrimSpace(password)

	if phone == "" || password == "" {
		return nil, "", errors.New("phone and password are required")
	}

	user, err := s.Users.GetByPhone(ctx, phone)
	if err != nil {
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

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) GenerateToken(userID int64, phone string) (string, error) {
	now := time.Now()

	claims := AuthClaims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "myapp-auth",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}

func (s *AuthService) ParseToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.JWTSecret, nil
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
