package service

type AuthService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}

type JWTService interface {
	GenerateToken(userID, email, role string) (string, error)
	ValidateToken(token string) (userID, email, role string, err error)
}
