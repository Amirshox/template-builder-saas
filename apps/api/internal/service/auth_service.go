package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"template-builder-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Repository
	// In a real app, this should be an env var
	jwtSecret []byte
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte("super-secret-key-change-me"),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (string, error) {
	// 1. Check if user exists
	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return "", errors.New("user already exists")
	}

	// 2. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 3. Create Org (1:1 for MVP)
	orgName := fmt.Sprintf("%s's Workspace", name)
	org, err := s.repo.CreateOrg(ctx, orgName)
	if err != nil {
		return "", fmt.Errorf("failed to create org: %w", err)
	}

	// 4. Create User
	user, err := s.repo.CreateUser(ctx, email, name, string(hashed))
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Create Membership (skipping formal DB generic link for MVP if not strictly enforced by RLS yet,
	// but ideally we should have it. Assuming repo doesn't have CreateMembership yet, we rely on knowing ID)

	// 5. Generate Token
	return s.GenerateToken(user.ID, org.ID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// For MVP, checking memberships or just picking headers?
	// Let's assume we find their primary org.
	// Since we don't have ListMemberships yet, we might need to HACK this for the MVP
	// or just require the client to switch orgs later.
	// HACK: For now, we need an OrgID to put in the token.
	// If we don't have a lookup, we'll encounter issues.
	// Let's add `GetOrgForUser` to repo or just query for now.

	// Simplification: We just registered them, so we know they have an org.
	// But Login needs to find it.
	// Let's add ListOrgsForUser to repo in next step.
	// For now, returning a token WITHOUT org, or defaulting to 0000 if not found (bad).
	// Let's rely on repo update in a moment.

	return s.GenerateToken(user.ID, uuid.Nil) // Placeholder OrgID, Middleware will complain or we fetch later?
}

func (s *AuthService) GenerateToken(userID, orgID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"org": orgID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
}
