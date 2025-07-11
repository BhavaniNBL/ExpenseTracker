package repository

import (
	"expensetracker/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	GetByID(id uuid.UUID) (*models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
	List(filters map[string]interface{}, from, to time.Time, limit, offset int) ([]models.Expense, error)
	Summary(userID string, from, to time.Time) (map[string]float64, error)
}

type expenseRepo struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepo{db}
}

func (r *expenseRepo) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepo) GetByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.First(&expense, "id = ?", id).Error
	return &expense, err
}

func (r *expenseRepo) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Expense{}, "id = ?", id).Error
}

func (r *expenseRepo) List(filters map[string]interface{}, from, to time.Time, limit, offset int) ([]models.Expense, error) {
	var expenses []models.Expense
	query := r.db.Where("timestamp BETWEEN ? AND ?", from, to)

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Limit(limit).Offset(offset).Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepo) Summary(userID string, from, to time.Time) (map[string]float64, error) {
	rows, err := r.db.
		Table("expenses").
		Select("category, SUM(amount) as total").
		Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, from, to).
		Group("category").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		summary[category] = total
	}
	return summary, nil
}
