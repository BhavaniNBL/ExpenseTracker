// package handlers

// import (
// 	"expensetracker/models"
// 	"expensetracker/services"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type ExpenseHandler struct {
// 	service services.ExpenseService
// }

// func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
// 	return &ExpenseHandler{service}
// }

// // CreateExpense godoc
// // @Summary Create a new expense
// // @Description Adds a new expense record for the authenticated user
// // @Tags expenses
// // @Accept json
// // @Produce json
// // @Param expense body models.Expense true "Expense object"
// // @Success 201 {object} models.Expense
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses [post]

// func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
// 	var exp models.Expense
// 	if err := c.ShouldBindJSON(&exp); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
// 		return
// 	}
// 	exp.Timestamp = time.Now()
// 	if uid, exists := c.Get("user_id"); exists {
// 		exp.UserID, _ = uuid.Parse(uid.(string))
// 	}
// 	// if err := h.service.Create(&exp); err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save expense"})
// 	// 	return
// 	// }
// 	if err := h.service.Create(&exp); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "could not save expense",
// 			"debug": err.Error(), // 👈 This will show the real error
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, exp)
// }

// // GetExpense godoc
// // @Summary Get an expense by ID
// // @Description Returns a single expense by its UUID
// // @Tags expenses
// // @Produce json
// // @Param id path string true "Expense ID"
// // @Success 200 {object} models.Expense
// // @Failure 404 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses/{id} [get]

// func (h *ExpenseHandler) GetExpense(c *gin.Context) {
// 	id := c.Param("id")
// 	uid, _ := uuid.Parse(id)
// 	exp, err := h.service.GetByID(uid)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, exp)
// }

// // UpdateExpense godoc
// // @Summary Update an expense
// // @Description Updates an existing expense by ID
// // @Tags expenses
// // @Accept json
// // @Produce json
// // @Param id path string true "Expense ID"
// // @Param expense body models.Expense true "Updated expense object"
// // @Success 200 {object} models.Expense
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses/{id} [put]

// func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
// 	id := c.Param("id")
// 	uid, _ := uuid.Parse(id)
// 	var exp models.Expense
// 	if err := c.ShouldBindJSON(&exp); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
// 		return
// 	}
// 	exp.ID = uid
// 	if err := h.service.Update(&exp); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, exp)
// }

// // DeleteExpense godoc
// // @Summary Delete an expense
// // @Description Deletes an expense by ID
// // @Tags expenses
// // @Produce json
// // @Param id path string true "Expense ID"
// // @Success 200 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses/{id} [delete]
// func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
// 	id := c.Param("id")
// 	uid, _ := uuid.Parse(id)
// 	if err := h.service.Delete(uid); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion failed"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
// }

// // ListExpenses godoc
// // @Summary List expenses
// // @Description Lists all expenses with optional filters
// // @Tags expenses
// // @Produce json
// // @Param category query string false "Category"
// // @Param currency query string false "Currency"
// // @Param from query string false "Start date (RFC3339)"
// // @Param to query string false "End date (RFC3339)"
// // @Param limit query int false "Limit"
// // @Param offset query int false "Offset"
// // @Success 200 {array} models.Expense
// // @Failure 500 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses [get]

// func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
// 	userID := c.Query("user_id")
// 	category := c.Query("category")
// 	currency := c.Query("currency")
// 	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", "2000-01-01T00:00:00Z"))
// 	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	exp, err := h.service.List(userID, category, currency, from, to, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch failed"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, exp)
// }

// // func (h *ExpenseHandler) SummaryExpenses(c *gin.Context) {

// // 	userID := c.Query("user_id")
// // 	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", "2000-01-01T00:00:00Z"))
// // 	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))
// // 	target := c.DefaultQuery("target_currency", "USD")

// // 	// summary, err := h.service.Summary(userID, from, to, target)
// // 	// if err != nil {
// // 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "summary failed"})
// // 	// 	return
// // 	// }
// // 	summary, err := h.service.Summary(userID, from, to, target)
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{
// // 			"error": "summary failed",
// // 			"debug": err.Error(), // 👈 show real error
// // 		})
// // 		return
// // 	}

// // 	c.JSON(http.StatusOK, summary)
// // }

// // SummaryExpenses godoc
// // @Summary Expense summary
// // @Description Provides category-wise expense summary for the user
// // @Tags expenses
// // @Produce json
// // @Param from query string false "Start date (RFC3339)"
// // @Param to query string false "End date (RFC3339)"
// // @Param target_currency query string false "Target currency (e.g. USD)"
// // @Success 200 {object} map[string]float64
// // @Failure 400 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Security BearerAuth
// // @Router /expenses/summary [get]
// func (h *ExpenseHandler) SummaryExpenses(c *gin.Context) {
// 	uidRaw, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
// 		return
// 	}
// 	userID := uidRaw.(string)

// 	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", "2000-01-01T00:00:00Z"))
// 	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))
// 	target := c.DefaultQuery("target_currency", "USD")

// 	summary, err := h.service.Summary(userID, from, to, target)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "summary failed", "debug": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, summary)
// }

package handlers

import (
	"expensetracker/models"
	"expensetracker/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExpenseHandler struct {
	service services.ExpenseService
}

func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service}
}

// CreateExpense godoc
// @Summary Create a new expense
// @Description Adds a new expense record for the authenticated user
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body models.Expense true "Expense object"
// @Success 201 {object} models.Expense
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var exp models.Expense
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	exp.Timestamp = time.Now()
	if uid, exists := c.Get("user_id"); exists {
		exp.UserID, _ = uuid.Parse(uid.(string))
	}
	if err := h.service.Create(&exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not save expense",
			"debug": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, exp)
}

// GetExpense godoc
// @Summary Get an expense by ID
// @Description Returns a single expense by its UUID
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} models.Expense
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses/{id} [get]
func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	id := c.Param("id")
	uid, _ := uuid.Parse(id)
	exp, err := h.service.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}
	c.JSON(http.StatusOK, exp)
}

// UpdateExpense godoc
// @Summary Update an expense
// @Description Updates an existing expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Param expense body models.Expense true "Updated expense object"
// @Success 200 {object} models.Expense
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	uid, _ := uuid.Parse(id)
	var exp models.Expense
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	exp.ID = uid
	if err := h.service.Update(&exp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, exp)
}

// DeleteExpense godoc
// @Summary Delete an expense
// @Description Deletes an expense by ID
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	uid, _ := uuid.Parse(id)
	if err := h.service.Delete(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deletion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ListExpenses godoc
// @Summary List expenses
// @Description Lists all expenses with optional filters
// @Tags expenses
// @Produce json
// @Param category query string false "Category"
// @Param currency query string false "Currency"
// @Param from query string false "Start date (RFC3339)"
// @Param to query string false "End date (RFC3339)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Expense
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses [get]
func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	userID := c.Query("user_id")
	category := c.Query("category")
	currency := c.Query("currency")
	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", "2000-01-01T00:00:00Z"))
	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	exp, err := h.service.List(userID, category, currency, from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch failed"})
		return
	}
	c.JSON(http.StatusOK, exp)
}

// SummaryExpenses godoc
// @Summary Expense summary
// @Description Provides category-wise expense summary for the user
// @Tags expenses
// @Produce json
// @Param from query string false "Start date (RFC3339)"
// @Param to query string false "End date (RFC3339)"
// @Param target_currency query string false "Target currency (e.g. USD)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/expenses/summary [get]
func (h *ExpenseHandler) SummaryExpenses(c *gin.Context) {
	uidRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}
	userID := uidRaw.(string)

	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", "2000-01-01T00:00:00Z"))
	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))
	target := c.DefaultQuery("target_currency", "USD")

	summary, err := h.service.Summary(userID, from, to, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "summary failed", "debug": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
