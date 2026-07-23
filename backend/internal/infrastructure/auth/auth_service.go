package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTServiceImpl struct {
	secret          []byte
	expirationHours int
}

func NewJWTService(cfg config.JWTConfig) *JWTServiceImpl {
	secret := cfg.Secret
	if secret == "" {
		env := os.Getenv("GO_ENV")
		if env == "production" || env == "staging" {
			panic("JWT_SECRET is required in production/staging")
		}
		secret = "dev-only-secret-not-for-production"
	}
	return &JWTServiceImpl{
		secret:          []byte(secret),
		expirationHours: cfg.ExpirationHours,
	}
}

func (s *JWTServiceImpl) GenerateToken(userID, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (userID, email, role string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return "", "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, claims.Email, claims.Role, nil
}

type BcryptAuthService struct {
	cost int
}

func NewBcryptAuthService(cost int) *BcryptAuthService {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptAuthService{cost: cost}
}

func (s *BcryptAuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	return string(bytes), err
}

func (s *BcryptAuthService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
