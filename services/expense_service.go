package services

import (
	"expensetracker/models"
	"expensetracker/repository"
	"time"

	"github.com/google/uuid"
)

type ExpenseService interface {
	Create(expense *models.Expense) error
	GetByID(id uuid.UUID) (*models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
	List(userID, category, currency string, from, to time.Time, limit, offset int) ([]models.Expense, error)
	Summary(userID string, from, to time.Time, targetCurrency string) (map[string]float64, error)
}

type expenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(repo repository.ExpenseRepository) ExpenseService {
	return &expenseService{repo}
}

func (s *expenseService) Create(expense *models.Expense) error {
	expense.ID = uuid.New()
	return s.repo.Create(expense)
}

func (s *expenseService) GetByID(id uuid.UUID) (*models.Expense, error) {
	return s.repo.GetByID(id)
}

func (s *expenseService) Update(expense *models.Expense) error {
	return s.repo.Update(expense)
}

func (s *expenseService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *expenseService) List(userID, category, currency string, from, to time.Time, limit, offset int) ([]models.Expense, error) {
	filters := make(map[string]interface{})
	if userID != "" {
		filters["user_id"] = userID
	}
	if category != "" {
		filters["category"] = category
	}
	if currency != "" {
		filters["currency"] = currency
	}
	return s.repo.List(filters, from, to, limit, offset)
}

func (s *expenseService) Summary(userID string, from, to time.Time, targetCurrency string) (map[string]float64, error) {
	raw, err := s.repo.Summary(userID, from, to)
	if err != nil {
		return nil, err
	}
	if targetCurrency != "" && targetCurrency != "USD" {
		// Mock conversion (e.g., assume all values are in USD, convert to target)
		return convertCurrency(raw, targetCurrency), nil
	}
	return raw, nil
}

func convertCurrency(input map[string]float64, target string) map[string]float64 {
	mockRates := map[string]float64{"EUR": 0.9, "INR": 82, "USD": 1.0}
	rate := mockRates[target]
	converted := make(map[string]float64)
	for cat, amt := range input {
		converted[cat] = amt * rate
	}
	return converted
}
