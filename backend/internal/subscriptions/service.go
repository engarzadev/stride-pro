package subscriptions

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Service provides subscription management operations.
// Currently stubbed for v2; checks are based on the user's subscription_tier.
type Service struct {
	db *database.DB
}

// NewService creates a subscription service.
func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

// GetCurrentPlan returns the plan for a user based on their subscription tier.
func (s *Service) GetCurrentPlan(userID uuid.UUID) (*Plan, error) {
	var tier string
	err := s.db.QueryRow("SELECT subscription_tier FROM users WHERE id = $1", userID).Scan(&tier)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying user tier: %w", err)
	}

	plan, ok := Plans[tier]
	if !ok {
		// Default to free plan if tier is unknown
		plan = Plans["free"]
	}
	return &plan, nil
}

// HasFeature checks whether a user's current plan includes a specific feature.
func (s *Service) HasFeature(userID uuid.UUID, feature string) (bool, error) {
	plan, err := s.GetCurrentPlan(userID)
	if err != nil {
		return false, err
	}

	for _, f := range plan.Features {
		if f == feature {
			return true, nil
		}
	}
	return false, nil
}

// GetSubscription returns the active subscription record for a user.
func (s *Service) GetSubscription(userID uuid.UUID) (*models.Subscription, error) {
	sub := &models.Subscription{}
	err := s.db.QueryRow(`
		SELECT id, user_id, plan, status, features, starts_at, ends_at, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(
		&sub.ID, &sub.UserID, &sub.Plan, &sub.Status, &sub.Features,
		&sub.StartsAt, &sub.EndsAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying subscription: %w", err)
	}
	return sub, nil
}

// ListPlans returns all available plans.
func (s *Service) ListPlans() []Plan {
	plans := make([]Plan, 0, len(Plans))
	for _, p := range Plans {
		plans = append(plans, p)
	}
	return plans
}

