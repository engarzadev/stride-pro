package sessions

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
	sessionDB, sessionMock := newTestDB(t)
	svc := NewService(
		NewRepository(sessionDB),
		subscriptions.NewService(subsDB),
	)
	return svc, subsMock, sessionMock
}

func validInput() CreateInput {
	return CreateInput{
		AppointmentID: uuid.New(),
		Type:          "massage",
		BodyZones:     []string{"back", "neck"},
		Notes:         "Patient was relaxed",
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
	svc, subsMock, sessionMock := newTestService(t)
	userID := uuid.New()
	input := validInput()

	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("base"))

	sessionMock.ExpectExec("INSERT INTO sessions").
		WithArgs(
			sqlmock.AnyArg(),    // id
			userID,              // user_id
			input.AppointmentID, // appointment_id
			"massage",           // type
			sqlmock.AnyArg(),    // body_zones (JSON)
			"Patient was relaxed", // notes
			sqlmock.AnyArg(),    // findings
			sqlmock.AnyArg(),    // recommendations
			sqlmock.AnyArg(),    // created_at
			sqlmock.AnyArg(),    // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, input)
	if err != nil {
		t.Fatalf("unexpected error on base plan: %v", err)
	}
}

func TestCreate_FeatureAllowed_EnterprisePlan(t *testing.T) {
	svc, subsMock, sessionMock := newTestService(t)
	userID := uuid.New()
	input := validInput()

	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("enterprise"))

	sessionMock.ExpectExec("INSERT INTO sessions").
		WithArgs(
			sqlmock.AnyArg(), userID, input.AppointmentID, "massage",
			sqlmock.AnyArg(), "Patient was relaxed",
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, input)
	if err != nil {
		t.Fatalf("unexpected error on enterprise plan: %v", err)
	}
}
