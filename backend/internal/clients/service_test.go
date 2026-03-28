package clients

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
	clientDB, clientMock := newTestDB(t)
	svc := NewService(
		NewRepository(clientDB),
		subscriptions.NewService(subsDB),
	)
	return svc, subsMock, clientMock
}

func validInput() CreateInput {
	return CreateInput{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}
}

func TestCreate_BelowLimit(t *testing.T) {
	svc, subsMock, clientMock := newTestService(t)
	userID := uuid.New()

	// Free plan: limit is 10; current count is 5
	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("free"))

	clientMock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	clientMock.ExpectExec("INSERT INTO clients").
		WithArgs(
			sqlmock.AnyArg(), // id
			userID,           // user_id
			"Jane",           // first_name
			"Smith",          // last_name
			"jane@example.com", // email
			sqlmock.AnyArg(), // phone
			sqlmock.AnyArg(), // address
			sqlmock.AnyArg(), // notes
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	client, err := svc.Create(userID, validInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.FirstName != "Jane" {
		t.Errorf("client.FirstName = %q, want \"Jane\"", client.FirstName)
	}
}

func TestCreate_AtLimit(t *testing.T) {
	svc, subsMock, clientMock := newTestService(t)
	userID := uuid.New()

	// Free plan: limit is 10; current count is already 10
	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("free"))

	clientMock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	_, err := svc.Create(userID, validInput())
	if !errors.Is(err, subscriptions.ErrLimitExceeded) {
		t.Errorf("expected ErrLimitExceeded, got: %v", err)
	}
}

func TestCreate_UnlimitedPlan(t *testing.T) {
	svc, subsMock, clientMock := newTestService(t)
	userID := uuid.New()

	// Base plan: no count check needed
	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("base"))

	clientMock.ExpectExec("INSERT INTO clients").
		WithArgs(
			sqlmock.AnyArg(),
			userID,
			"Jane",
			"Smith",
			"jane@example.com",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, validInput())
	if err != nil {
		t.Fatalf("unexpected error on unlimited plan: %v", err)
	}
}

func TestCreate_ExactlyAtLimit_OneBelow(t *testing.T) {
	svc, subsMock, clientMock := newTestService(t)
	userID := uuid.New()

	// Count is 9: one slot left
	subsMock.ExpectQuery("SELECT subscription_tier FROM users").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"subscription_tier"}).AddRow("free"))

	clientMock.ExpectQuery("SELECT COUNT").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(9))

	clientMock.ExpectExec("INSERT INTO clients").
		WithArgs(
			sqlmock.AnyArg(), userID, "Jane", "Smith", "jane@example.com",
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := svc.Create(userID, validInput())
	if err != nil {
		t.Fatalf("unexpected error at count=9 (limit=10): %v", err)
	}
}
