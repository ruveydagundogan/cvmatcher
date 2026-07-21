package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/config"
)

type JWTServiceImpl struct {
	secret          []byte
	expirationHours int
}

func NewJWTService(cfg config.JWTConfig) *JWTServiceImpl {
	return &JWTServiceImpl{
		secret:          []byte(cfg.Secret),
		expirationHours: cfg.ExpirationHours,
	}
}

func (s *JWTServiceImpl) GenerateToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(s.expirationHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (userID, email, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", "", fmt.Errorf("invalid token claims")
	}

	userID, _ = claims["user_id"].(string)
	email, _ = claims["email"].(string)
	role, _ = claims["role"].(string)

	return userID, email, role, nil
}

type BcryptAuthService struct{}

func NewBcryptAuthService() *BcryptAuthService {
	return &BcryptAuthService{}
}

func (s *BcryptAuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *BcryptAuthService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
