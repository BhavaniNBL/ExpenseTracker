package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	handlers "expensetracker/handlers"
	"expensetracker/models"
	"expensetracker/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockExpenseService implements ExpenseService for testing
type mockExpenseService struct {
	CreateFn  func(*models.Expense) error
	GetByIDFn func(uuid.UUID) (*models.Expense, error)
	UpdateFn  func(*models.Expense) error
	DeleteFn  func(uuid.UUID) error
	ListFn    func(string, string, string, time.Time, time.Time, int, int) ([]models.Expense, error)
	SummaryFn func(string, time.Time, time.Time, string) (map[string]float64, error)
}

func (m *mockExpenseService) Create(e *models.Expense) error {
	return m.CreateFn(e)
}
func (m *mockExpenseService) GetByID(id uuid.UUID) (*models.Expense, error) {
	return m.GetByIDFn(id)
}
func (m *mockExpenseService) Update(e *models.Expense) error {
	return m.UpdateFn(e)
}
func (m *mockExpenseService) Delete(id uuid.UUID) error {
	return m.DeleteFn(id)
}
func (m *mockExpenseService) List(userID, category, currency string, from, to time.Time, limit, offset int) ([]models.Expense, error) {
	return m.ListFn(userID, category, currency, from, to, limit, offset)
}
func (m *mockExpenseService) Summary(userID string, from, to time.Time, targetCurrency string) (map[string]float64, error) {
	return m.SummaryFn(userID, from, to, targetCurrency)
}

func setupRouter(service services.ExpenseService) (*gin.Engine, *handlers.ExpenseHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewExpenseHandler(service)
	return r, handler
}

func TestCreateExpense_Success(t *testing.T) {
	mockSvc := &mockExpenseService{
		CreateFn: func(e *models.Expense) error {
			e.ID = uuid.New()
			return nil
		},
	}
	_, h := setupRouter(mockSvc)

	exp := models.Expense{Amount: 123.45, Description: "Test expense"}
	body, _ := json.Marshal(exp)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uuid.New().String())

	h.CreateExpense(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Expense
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, exp.Amount, resp.Amount)
	assert.NotEmpty(t, resp.ID)
}

func TestCreateExpense_InvalidJSON(t *testing.T) {
	mockSvc := &mockExpenseService{}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/expenses", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateExpense(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid input")
}

func TestGetExpense_Success(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{
		GetByIDFn: func(uid uuid.UUID) (*models.Expense, error) {
			return &models.Expense{ID: uid, Amount: 42}, nil
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/expenses/"+id.String(), nil)

	h.GetExpense(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Expense
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
}

func TestGetExpense_NotFound(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{
		GetByIDFn: func(uid uuid.UUID) (*models.Expense, error) {
			return nil, errors.New("not found")
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/expenses/"+id.String(), nil)

	h.GetExpense(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "expense not found")
}

func TestUpdateExpense_Success(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{
		UpdateFn: func(e *models.Expense) error { return nil },
	}
	_, h := setupRouter(mockSvc)

	exp := models.Expense{Amount: 78.9, Description: "Updated expense"}
	body, _ := json.Marshal(exp)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/expenses/"+id.String(), bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.UpdateExpense(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Expense
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, id, resp.ID)
}

func TestUpdateExpense_InvalidJSON(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/expenses/"+id.String(), bytes.NewBufferString("bad json"))
	c.Request.Header.Set("Content-Type", "application/json")

	h.UpdateExpense(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid input")
}

func TestDeleteExpense_Success(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{
		DeleteFn: func(uid uuid.UUID) error { return nil },
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/api/v1/expenses/"+id.String(), nil)

	h.DeleteExpense(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted")
}

func TestDeleteExpense_Error(t *testing.T) {
	id := uuid.New()
	mockSvc := &mockExpenseService{
		DeleteFn: func(uid uuid.UUID) error { return errors.New("fail") },
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/api/v1/expenses/"+id.String(), nil)

	h.DeleteExpense(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "deletion failed")
}

func TestListExpenses_Success(t *testing.T) {
	mockSvc := &mockExpenseService{
		ListFn: func(userID, category, currency string, from, to time.Time, limit, offset int) ([]models.Expense, error) {
			return []models.Expense{
				{ID: uuid.New(), Amount: 100},
				{ID: uuid.New(), Amount: 50},
			}, nil
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/expenses?limit=2&offset=0", nil)
	c.Request = req

	h.ListExpenses(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []models.Expense
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestListExpenses_Error(t *testing.T) {
	mockSvc := &mockExpenseService{
		ListFn: func(userID, category, currency string, from, to time.Time, limit, offset int) ([]models.Expense, error) {
			return nil, errors.New("fail")
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/expenses", nil)
	c.Request = req

	h.ListExpenses(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "fetch failed")
}

func TestSummaryExpenses_Success(t *testing.T) {
	mockSvc := &mockExpenseService{
		SummaryFn: func(userID string, from, to time.Time, targetCurrency string) (map[string]float64, error) {
			return map[string]float64{"food": 200}, nil
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/expenses/summary", nil)
	c.Set("user_id", "user123")

	h.SummaryExpenses(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]float64
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 200.0, resp["food"])
}

func TestSummaryExpenses_MissingUserID(t *testing.T) {
	mockSvc := &mockExpenseService{}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/expenses/summary", nil)
	// user_id not set in context

	h.SummaryExpenses(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "user_id not found in context")
}

func TestSummaryExpenses_Error(t *testing.T) {
	mockSvc := &mockExpenseService{
		SummaryFn: func(userID string, from, to time.Time, targetCurrency string) (map[string]float64, error) {
			return nil, errors.New("fail")
		},
	}
	_, h := setupRouter(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/expenses/summary", nil)
	c.Set("user_id", "user123")

	h.SummaryExpenses(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "summary failed")
}
