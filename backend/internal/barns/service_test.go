package barns

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/subscriptions"
)

func newTestDB(t *testing.T) (*database.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("creating sqlmock: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	return &database.DB{DB: sqlDB}, mock
}

func newTestService(t *testing.T) (*Service, sqlmock.Sqlmock, sqlmock.Sqlmock) {
	t.Helper()
	subsDB, subsMock := newTestDB(t)
	barnDB, barnMock := newTestDB(t)
	svc := NewService(
		NewRepository(barnDB),
		subscriptions.NewService(subsDB),
	)
	return svc, subsMock, barnMock
}

func validInput() CreateInput {
	return CreateInput{
		Name:  "Sunny Stables",
		Email: "barn@example.com",
	}
}

func TestCreate_FeatureGated_FreePlan(t *testing.T) {
	svc, subsMock, _ := newTestService(t)
	userID := uuid.New()

	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("free"))

	_, err := svc.Create(userID, validInput())
	if !errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
		t.Errorf("expected ErrFeatureNotAvailable for free plan, got: %v", err)
	}
}

func TestCreate_FeatureAllowed_BasePlan(t *testing.T) {
	svc, subsMock, barnMock := newTestService(t)
	userID := uuid.New()

	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("base"))

	barnMock.ExpectExec("INSERT INTO barns").
		WithArgs(
			sqlmock.AnyArg(), // id
			userID,           // user_id
			"Sunny Stables",  // name
			sqlmock.AnyArg(), // contact_name
			sqlmock.AnyArg(), // address
			sqlmock.AnyArg(), // phone
			"barn@example.com", // email
			sqlmock.AnyArg(), // notes
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, validInput())
	if err != nil {
		t.Fatalf("unexpected error on base plan: %v", err)
	}
}

func TestCreate_FeatureAllowed_TrainerAddonPlan(t *testing.T) {
	svc, subsMock, barnMock := newTestService(t)
	userID := uuid.New()

	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("trainer_addon"))

	barnMock.ExpectExec("INSERT INTO barns").
		WithArgs(
			sqlmock.AnyArg(), userID, "Sunny Stables",
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"barn@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, validInput())
	if err != nil {
		t.Fatalf("unexpected error on trainer_addon plan: %v", err)
	}
}
