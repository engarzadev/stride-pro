// Package auth handles user authentication including registration, login, and JWT management.
package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

var (
	ErrEmailTaken       = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound     = errors.New("user not found")
)

// Service handles authentication business logic.
type Service struct {
	db        *database.DB
	jwtSecret []byte
}

// NewService creates an auth service with the given database and JWT secret.
func NewService(db *database.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// RegisterInput holds the data needed to create a new user.
type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginInput holds login credentials.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// Register creates a new user with a hashed password.
func (s *Service) Register(input RegisterInput) (*models.User, *TokenPair, error) {
	// Check if email already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", input.Email).Scan(&exists)
	if err != nil {
		return nil, nil, fmt.Errorf("checking email existence: %w", err)
	}
	if exists {
		return nil, nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		ID:               uuid.New(),
		Email:            input.Email,
		PasswordHash:     string(hash),
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		Role:             "user",
		SubscriptionTier: "free",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = s.db.Exec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, subscription_tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.SubscriptionTier, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("inserting user: %w", err)
	}

	tokens, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Login authenticates a user and returns a token pair.
func (s *Service) Login(input LoginInput) (*models.User, *TokenPair, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, role, subscription_tier, created_at, updated_at
		FROM users WHERE email = $1`,
		input.Email,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SubscriptionTier, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, nil, fmt.Errorf("querying user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// ValidateToken parses and validates a JWT, returning the user ID.
func (s *Service) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("missing subject claim")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing user ID from token: %w", err)
	}

	return userID, nil
}

// RefreshToken generates a new token pair from a valid refresh token.
func (s *Service) RefreshToken(refreshToken string) (*TokenPair, error) {
	userID, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}
	return s.generateTokenPair(userID)
}

// GetUserByID looks up a user by their ID.
func (s *Service) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, role, subscription_tier, created_at, updated_at
		FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SubscriptionTier, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by ID: %w", err)
	}
	return user, nil
}

func (s *Service) generateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	accessClaims := jwt.MapClaims{
		"sub":  userID.String(),
		"type": "access",
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub":  userID.String(),
		"type": "refresh",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshTokenJWT.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
