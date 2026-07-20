package service

import (
	"context"
	"hash"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
)

type AuthService struct {
	users repository.UserRepository
	tokens *jwt.Manager
}

func NewAuthService(users repository.UserRepository, tokens *jwt.Manager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

type AuthResult struct {
	Token string `json:"token"`
	User *model.User `json:"user"`
}

type RegisterInput struct {
	Username string
	Email string
	Password string
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	// 1. Validation
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
	if in.Username == "" || in.Email == "" || in.Password == "" {
		return nil, apperror.New(apperror.CodeInvalidArgument, "username, email and password are required")
	}
	if len(in.Password) < 8 {
		return nil, apperror.New(apperror.CodeInvalidArgument, "password must be at least 8 characters")
	}

	// 2. Hash password
	const Cost = 12
	b, err := bcrypt.GenerateFromPassword([]byte(in.Password), Cost)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "hashing password", err)
	}

	// 3. Create user
	u := &model.User{
		Username: in.Username,
		Email: in.Email,
		PasswordHash: string(b),
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	// 4. Issue a JWT
	token, err := s.tokens.Generate(u.ID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "issuing token", err)
	}

	// TODO: 5. Audit log

	return &AuthResult{Token: token, User: u}, nil
}

