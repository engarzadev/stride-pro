// Package auth handles user authentication including registration, login, and JWT management.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	ErrEmailTaken              = errors.New("email already registered")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrUserNotFound            = errors.New("user not found")
	ErrTokenRevoked            = errors.New("token has been revoked")
	ErrResetTokenInvalid       = errors.New("reset token is invalid or expired")
)

// PasswordResetter is satisfied by notifications.EmailSender; defined here to
// avoid an import cycle between auth and notifications.
type PasswordResetter interface {
	SendHTML(recipient, subject, htmlBody string) error
}

// Service handles authentication business logic.
type Service struct {
	db          *database.DB
	jwtSecret   []byte
	emailSender PasswordResetter
	appBaseURL  string // e.g. "https://app.stridepro.com" — used to build reset links
}

// NewService creates an auth service with the given database and JWT secret.
func NewService(db *database.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// SetEmailSender configures the email sender used for password reset emails.
// Call this from main after constructing the service.
func (s *Service) SetEmailSender(sender PasswordResetter, appBaseURL string) {
	s.emailSender = sender
	s.appBaseURL = appBaseURL
}

// RegisterInput holds the data needed to create a new user.
type RegisterInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	AccountType string `json:"account_type"`
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

	role := "professional"
	if input.AccountType == "owner" {
		role = "owner"
	}

	user := &models.User{
		ID:               uuid.New(),
		Email:            input.Email,
		PasswordHash:     string(hash),
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		Role:             role,
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

// ValidateToken parses and validates a JWT access token, returning the user ID.
// It rejects refresh tokens, expired tokens, and revoked tokens.
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

	// Reject refresh tokens being used as access tokens
	if tokenType, _ := claims["type"].(string); tokenType != "access" {
		return uuid.Nil, errors.New("invalid token type")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("missing subject claim")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing user ID from token: %w", err)
	}

	// Check revocation list
	jti, _ := claims["jti"].(string)
	if jti != "" {
		revoked, err := s.isTokenRevoked(jti)
		if err != nil {
			return uuid.Nil, fmt.Errorf("checking token revocation: %w", err)
		}
		if revoked {
			return uuid.Nil, ErrTokenRevoked
		}
	}

	return userID, nil
}

// RevokeToken adds the given token's JTI to the revoked tokens table, invalidating it
// immediately even before its natural expiry.
func (s *Service) RevokeToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	jti, _ := claims["jti"].(string)
	if jti == "" {
		// Token predates JTI support — nothing to revoke
		return nil
	}

	// Determine the token's expiry so the row can be cleaned up later
	exp, _ := claims["exp"].(float64)
	expiresAt := time.Unix(int64(exp), 0)

	_, err = s.db.Exec(
		`INSERT INTO revoked_tokens (jti, expires_at) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`,
		jti, expiresAt,
	)
	return err
}

// RefreshToken generates a new token pair from a valid refresh token.
func (s *Service) RefreshToken(refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Only accept refresh tokens here
	if tokenType, _ := claims["type"].(string); tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("missing subject claim")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID from token: %w", err)
	}

	return s.generateTokenPair(userID)
}

// UpdateProfileInput holds the fields a user can update on their profile.
type UpdateProfileInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// ChangePasswordInput holds the current and new passwords for a password change.
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UpdateProfile updates first_name, last_name, and email for the given user.
func (s *Service) UpdateProfile(userID uuid.UUID, input UpdateProfileInput) (*models.User, error) {
	// If email is changing, make sure it isn't already taken by another account
	var existingID uuid.UUID
	err := s.db.QueryRow(`SELECT id FROM users WHERE email = $1`, input.Email).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("checking email: %w", err)
	}
	if err == nil && existingID != userID {
		return nil, ErrEmailTaken
	}

	user := &models.User{}
	err = s.db.QueryRow(`
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, email, password_hash, first_name, last_name, role, subscription_tier, created_at, updated_at`,
		input.FirstName, input.LastName, input.Email, userID,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Role,
		&user.SubscriptionTier, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("updating profile: %w", err)
	}
	return user, nil
}

// ChangePassword verifies the current password then replaces it with the new one.
func (s *Service) ChangePassword(userID uuid.UUID, input ChangePasswordInput) error {
	var hash string
	err := s.db.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("fetching user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if _, err = s.db.Exec(
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		string(newHash), userID,
	); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}
	return nil
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
		"jti":  uuid.New().String(),
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
		"jti":  uuid.New().String(),
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

func (s *Service) isTokenRevoked(jti string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti = $1 AND expires_at > NOW())`,
		jti,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ForgotPassword looks up the user by email, creates a single-use reset token
// valid for 1 hour, and sends a reset-link email. If no account matches the
// email we return nil (no enumeration — the caller always returns 200).
func (s *Service) ForgotPassword(email string) error {
	var userID uuid.UUID
	var firstName string
	err := s.db.QueryRow(
		`SELECT id, first_name FROM users WHERE email = $1`, email,
	).Scan(&userID, &firstName)
	if errors.Is(err, sql.ErrNoRows) {
		// Do not reveal whether the email is registered.
		return nil
	}
	if err != nil {
		return fmt.Errorf("querying user by email: %w", err)
	}

	// Generate a cryptographically random 32-byte token.
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Errorf("generating reset token: %w", err)
	}
	rawHex := hex.EncodeToString(raw) // sent to the user in the URL
	hash := sha256.Sum256([]byte(rawHex))
	tokenHash := hex.EncodeToString(hash[:]) // stored in the DB

	expiresAt := time.Now().Add(1 * time.Hour)
	_, err = s.db.Exec(
		`INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("storing reset token: %w", err)
	}

	if s.emailSender != nil {
		resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.appBaseURL, rawHex)
		subject, body := passwordResetEmail(firstName, resetURL)
		// Best-effort — log but don't fail; the token is already stored.
		if err := s.emailSender.SendHTML(email, subject, body); err != nil {
			return fmt.Errorf("sending reset email: %w", err)
		}
	}

	return nil
}

// ResetPassword validates the raw token, checks it hasn't been used or expired,
// then updates the user's password and marks the token as used.
func (s *Service) ResetPassword(rawToken, newPassword string) error {
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	var tokenID uuid.UUID
	var userID uuid.UUID
	var usedAt sql.NullTime
	err := s.db.QueryRow(
		`SELECT id, user_id, used_at FROM password_reset_tokens
		 WHERE token_hash = $1 AND expires_at > NOW()`,
		tokenHash,
	).Scan(&tokenID, &userID, &usedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrResetTokenInvalid
	}
	if err != nil {
		return fmt.Errorf("looking up reset token: %w", err)
	}
	if usedAt.Valid {
		return ErrResetTokenInvalid
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing new password: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err = tx.Exec(
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		string(newHash), userID,
	); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	if _, err = tx.Exec(
		`UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1`,
		tokenID,
	); err != nil {
		return fmt.Errorf("marking token used: %w", err)
	}

	return tx.Commit()
}

// passwordResetEmail returns the subject and HTML body for a password reset email.
func passwordResetEmail(firstName, resetURL string) (subject, body string) {
	subject = "Reset your Stride Pro password"
	body = fmt.Sprintf(`<p>Hi %s,</p>
<p>We received a request to reset your Stride Pro password. Click the button below to choose a new one. This link expires in 1 hour.</p>
<p><a href="%s" style="background:#3b6255;color:#fff;padding:12px 24px;border-radius:6px;text-decoration:none;display:inline-block;">Reset Password</a></p>
<p>If you didn't request this, you can safely ignore this email — your password won't change.</p>
<p>— The Stride Pro Team</p>`, firstName, resetURL)
	return
}
