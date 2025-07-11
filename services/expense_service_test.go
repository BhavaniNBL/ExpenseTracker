package services_test

import (
	"errors"
	"expensetracker/models"
	"expensetracker/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mock repository implementing ExpenseRepository
type mockRepo struct {
	CreateFn  func(*models.Expense) error
	GetFn     func(uuid.UUID) (*models.Expense, error)
	UpdateFn  func(*models.Expense) error
	DeleteFn  func(uuid.UUID) error
	ListFn    func(map[string]interface{}, time.Time, time.Time, int, int) ([]models.Expense, error)
	SummaryFn func(string, time.Time, time.Time) (map[string]float64, error)
}

func (m *mockRepo) Create(e *models.Expense) error {
	return m.CreateFn(e)
}
func (m *mockRepo) GetByID(id uuid.UUID) (*models.Expense, error) {
	return m.GetFn(id)
}
func (m *mockRepo) Update(e *models.Expense) error {
	return m.UpdateFn(e)
}
func (m *mockRepo) Delete(id uuid.UUID) error {
	return m.DeleteFn(id)
}
func (m *mockRepo) List(filters map[string]interface{}, from, to time.Time, limit, offset int) ([]models.Expense, error) {
	return m.ListFn(filters, from, to, limit, offset)
}
func (m *mockRepo) Summary(userID string, from, to time.Time) (map[string]float64, error) {
	return m.SummaryFn(userID, from, to)
}

// Tests for Create and GetByID already exist, add Update test
func TestUpdateSuccess(t *testing.T) {
	mock := &mockRepo{
		UpdateFn: func(e *models.Expense) error {
			return nil
		},
	}
	svc := services.NewExpenseService(mock)
	err := svc.Update(&models.Expense{})
	assert.NoError(t, err)
}

// Test Delete method
func TestDeleteSuccess(t *testing.T) {
	mock := &mockRepo{
		DeleteFn: func(id uuid.UUID) error {
			return nil
		},
	}
	svc := services.NewExpenseService(mock)
	err := svc.Delete(uuid.New())
	assert.NoError(t, err)
}

// Test List method with no filters
func TestListNoFilters(t *testing.T) {
	mock := &mockRepo{
		ListFn: func(filters map[string]interface{}, from, to time.Time, limit, offset int) ([]models.Expense, error) {
			assert.Empty(t, filters) // no filters applied
			return []models.Expense{{Amount: 100}}, nil
		},
	}
	svc := services.NewExpenseService(mock)
	results, err := svc.List("", "", "", time.Time{}, time.Time{}, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

// Test List method with filters
func TestListWithFilters(t *testing.T) {
	mock := &mockRepo{
		ListFn: func(filters map[string]interface{}, from, to time.Time, limit, offset int) ([]models.Expense, error) {
			assert.Equal(t, "user1", filters["user_id"])
			assert.Equal(t, "food", filters["category"])
			assert.Equal(t, "USD", filters["currency"])
			return []models.Expense{{Amount: 200}}, nil
		},
	}
	svc := services.NewExpenseService(mock)
	results, err := svc.List("user1", "food", "USD", time.Time{}, time.Time{}, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

// Test Summary method with error from repo
func TestSummaryError(t *testing.T) {
	mock := &mockRepo{
		SummaryFn: func(userID string, from, to time.Time) (map[string]float64, error) {
			return nil, errors.New("some error")
		},
	}
	svc := services.NewExpenseService(mock)
	_, err := svc.Summary("user1", time.Time{}, time.Time{}, "")
	assert.Error(t, err)
}

// Test Summary method with no conversion (targetCurrency empty or USD)
func TestSummaryNoConversion(t *testing.T) {
	mock := &mockRepo{
		SummaryFn: func(userID string, from, to time.Time) (map[string]float64, error) {
			return map[string]float64{"food": 50}, nil
		},
	}
	svc := services.NewExpenseService(mock)
	// targetCurrency empty
	result, err := svc.Summary("user1", time.Time{}, time.Time{}, "")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, result["food"])

	// targetCurrency USD
	result, err = svc.Summary("user1", time.Time{}, time.Time{}, "USD")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, result["food"])
}

// Test Summary method with currency conversion
func TestSummaryWithConversion(t *testing.T) {
	mock := &mockRepo{
		SummaryFn: func(userID string, from, to time.Time) (map[string]float64, error) {
			return map[string]float64{"food": 100}, nil
		},
	}
	svc := services.NewExpenseService(mock)
	result, err := svc.Summary("user1", time.Time{}, time.Time{}, "EUR")
	assert.NoError(t, err)
	assert.InDelta(t, 90.0, result["food"], 0.001)

	result, err = svc.Summary("user1", time.Time{}, time.Time{}, "INR")
	assert.NoError(t, err)
	assert.InDelta(t, 8200.0, result["food"], 0.001)
}

func TestSummaryWithoutConversion(t *testing.T) {
	mock := &mockRepo{
		SummaryFn: func(userID string, from, to time.Time) (map[string]float64, error) {
			return map[string]float64{"food": 100}, nil
		},
	}
	svc := services.NewExpenseService(mock)

	// Case 1: targetCurrency is empty string
	result, err := svc.Summary("user1", time.Time{}, time.Time{}, "")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, result["food"])

	// Case 2: targetCurrency is "USD"
	result, err = svc.Summary("user1", time.Time{}, time.Time{}, "USD")
	assert.NoError(t, err)
	assert.Equal(t, 100.0, result["food"])
}
