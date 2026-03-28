package subscriptions

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
)

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("creating sqlmock: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	return NewService(&database.DB{DB: sqlDB}), mock
}

func TestGetCurrentPlan_KnownTier(t *testing.T) {
	cases := []struct {
		tier     string
		wantID   string
		wantName string
	}{
		{"free", "free", "Free"},
		{"base", "base", "Base"},
		{"trainer_addon", "trainer_addon", "Trainer Add-on"},
		{"enterprise", "enterprise", "Enterprise"},
	}

	for _, tc := range cases {
		t.Run(tc.tier, func(t *testing.T) {
			svc, mock := newTestService(t)
			userID := uuid.New()

			mock.ExpectQuery("SELECT subscription_tier FROM users").
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow(tc.tier))

			plan, err := svc.GetCurrentPlan(userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if plan.ID != tc.wantID {
				t.Errorf("plan.ID = %q, want %q", plan.ID, tc.wantID)
			}
			if plan.Name != tc.wantName {
				t.Errorf("plan.Name = %q, want %q", plan.Name, tc.wantName)
			}
		})
	}
}

func TestGetCurrentPlan_UnknownTierDefaultsFree(t *testing.T) {
	svc, mock := newTestService(t)
	userID := uuid.New()

	mock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("legacy_plan"))

	plan, err := svc.GetCurrentPlan(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan.ID != "free" {
		t.Errorf("plan.ID = %q, want \"free\" for unknown tier", plan.ID)
	}
}

func TestGetClientLimit(t *testing.T) {
	cases := []struct {
		tier      string
		wantLimit int
	}{
		{"free", 10},
		{"base", -1},
		{"trainer_addon", -1},
		{"enterprise", -1},
	}

	for _, tc := range cases {
		t.Run(tc.tier, func(t *testing.T) {
			svc, mock := newTestService(t)
			userID := uuid.New()

			mock.ExpectQuery("SELECT subscription_tier FROM users").
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow(tc.tier))

			limit, err := svc.GetClientLimit(userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if limit != tc.wantLimit {
				t.Errorf("GetClientLimit(%q) = %d, want %d", tc.tier, limit, tc.wantLimit)
			}
		})
	}
}

func TestGetHorseLimit(t *testing.T) {
	cases := []struct {
		tier      string
		wantLimit int
	}{
		{"free", 20},
		{"base", -1},
		{"trainer_addon", -1},
		{"enterprise", -1},
	}

	for _, tc := range cases {
		t.Run(tc.tier, func(t *testing.T) {
			svc, mock := newTestService(t)
			userID := uuid.New()

			mock.ExpectQuery("SELECT subscription_tier FROM users").
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow(tc.tier))

			limit, err := svc.GetHorseLimit(userID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if limit != tc.wantLimit {
				t.Errorf("GetHorseLimit(%q) = %d, want %d", tc.tier, limit, tc.wantLimit)
			}
		})
	}
}

func TestHasFeature(t *testing.T) {
	cases := []struct {
		tier    string
		feature string
		want    bool
	}{
		{"free", "clients_max_10", true},
		{"free", "barn_management", false},
		{"free", "session_notes", false},
		{"base", "barn_management", true},
		{"base", "session_notes", true},
		{"base", "sms_notifications", false},
		{"trainer_addon", "sms_notifications", true},
		{"enterprise", "api_access", true},
		{"free", "api_access", false},
	}

	for _, tc := range cases {
		t.Run(tc.tier+"/"+tc.feature, func(t *testing.T) {
			svc, mock := newTestService(t)
			userID := uuid.New()

			mock.ExpectQuery("SELECT subscription_tier FROM users").
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow(tc.tier))

			got, err := svc.HasFeature(userID, tc.feature)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("HasFeature(%q, %q) = %v, want %v", tc.tier, tc.feature, got, tc.want)
			}
		})
	}
}

func TestRequireFeature_Available(t *testing.T) {
	svc, mock := newTestService(t)
	userID := uuid.New()

	mock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("base"))

	if err := svc.RequireFeature(userID, "barn_management"); err != nil {
		t.Errorf("expected nil error for available feature, got: %v", err)
	}
}

func TestRequireFeature_NotAvailable(t *testing.T) {
	features := []string{"barn_management", "session_notes", "sms_notifications"}

	for _, feature := range features {
		t.Run(feature, func(t *testing.T) {
			svc, mock := newTestService(t)
			userID := uuid.New()

			mock.ExpectQuery("SELECT subscription_tier FROM users").
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("free"))

			err := svc.RequireFeature(userID, feature)
			if err != ErrFeatureNotAvailable {
				t.Errorf("RequireFeature(%q) on free plan: got %v, want ErrFeatureNotAvailable", feature, err)
			}
		})
	}
}
