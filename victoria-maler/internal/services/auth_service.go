package services

import (
	"database/sql"
	"errors"
	"time"

	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/utils"
)

var (
	ErrEmailTaken  = errors.New("email already registered")
	ErrBadLogin    = errors.New("invalid email or password")
	ErrInvalidRole = errors.New("invalid role")
)

type AuthService struct {
	users  repositories.UserRepository
	secret string
}

func NewAuthService(users repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{users: users, secret: cfg.JWTSecret}
}

func (s *AuthService) Register(name, email, password, role string) (*models.User, error) {
	if role == "" {
		role = models.RoleDonor
	}
	if role != models.RoleDonor && role != models.RoleNGO && role != models.RoleAdmin {
		return nil, ErrInvalidRole
	}

	if _, err := s.users.FindByEmail(email); err == nil {
		return nil, ErrEmailTaken
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, ErrBadLogin
	}
	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, ErrBadLogin
	}
	return user, nil
}

func (s *AuthService) UserByID(id int64) (*models.User, error) {
	return s.users.FindByID(id)
}

func (s *AuthService) Token(user *models.User) (string, error) {
	return utils.GenerateToken(s.secret, user.ID, user.Email, user.Role)
}

func (s *AuthService) ValidateToken(token string) (*utils.Claims, error) {
	return utils.ParseToken(s.secret, token)
}
