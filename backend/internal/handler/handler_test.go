package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ppablomunoz/ownpocket/backend/internal/config"
	"github.com/ppablomunoz/ownpocket/backend/internal/middleware"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/ppablomunoz/ownpocket/backend/internal/service"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	require.NoError(t, err)

	db.Exec("PRAGMA foreign_keys = ON;")

	err = db.AutoMigrate(
		&model.User{},
		&model.Currency{},
		&model.Account{},
		&model.Category{},
		&model.Transaction{},
		&model.Budget{},
		&model.Tag{},
		&model.TransactionTag{},
	)
	require.NoError(t, err)

	seedTestCurrencies(t, db)
	return db
}

func seedTestCurrencies(t *testing.T, db *gorm.DB) {
	t.Helper()
	var count int64
	db.Model(&model.Currency{}).Count(&count)
	if count > 0 {
		return
	}
	currencies := []model.Currency{
		{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
		{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
	}
	err := db.Create(&currencies).Error
	require.NoError(t, err)
}

func setupTestRouter(t *testing.T) (*gin.Engine, *Handler, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	cfg := &config.Config{JWTSecret: "test-secret"}
	h := NewHandler(db, cfg)
	svc := service.NewService(db)

	// Recreate handler with all dependencies
	h.service = svc
	h.db = db
	h.cfg = cfg

	r := gin.New()
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")

	v1.GET("/health", h.Health)
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)

	protected := v1.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	protected.GET("/accounts", h.GetAccounts)
	protected.POST("/accounts", h.CreateAccount)
	protected.GET("/accounts/:id", h.GetAccount)
	protected.PUT("/accounts/:id", h.UpdateAccount)
	protected.DELETE("/accounts/:id", h.DeleteAccount)

	protected.GET("/transactions", h.GetTransactions)
	protected.POST("/transactions", h.CreateTransaction)
	protected.GET("/transactions/:id", h.GetTransaction)
	protected.PUT("/transactions/:id", h.UpdateTransaction)
	protected.DELETE("/transactions/:id", h.DeleteTransaction)

	protected.GET("/categories", h.GetCategories)
	protected.POST("/categories", h.CreateCategory)

	protected.GET("/budgets", h.GetBudgets)
	protected.POST("/budgets", h.CreateBudget)

	protected.GET("/dashboard/summary", h.GetDashboardSummary)

	return r, h, db
}

func setupTestRouterWithAuth(t *testing.T) (*gin.Engine, *Handler, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	cfg := &config.Config{JWTSecret: "test-secret"}
	h := NewHandler(db, cfg)
	svc := service.NewService(db)
	h.service = svc

	r := gin.New()
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")

	v1.GET("/health", h.Health)
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret)))

	protected.GET("/accounts", h.GetAccounts)
	protected.POST("/accounts", h.CreateAccount)
	protected.GET("/accounts/:id", h.GetAccount)
	protected.PUT("/accounts/:id", h.UpdateAccount)
	protected.DELETE("/accounts/:id", h.DeleteAccount)

	protected.GET("/transactions", h.GetTransactions)
	protected.POST("/transactions", h.CreateTransaction)
	protected.GET("/transactions/:id", h.GetTransaction)
	protected.PUT("/transactions/:id", h.UpdateTransaction)
	protected.DELETE("/transactions/:id", h.DeleteTransaction)

	protected.GET("/categories", h.GetCategories)
	protected.POST("/categories", h.CreateCategory)

	protected.GET("/budgets", h.GetBudgets)
	protected.POST("/budgets", h.CreateBudget)

	protected.GET("/dashboard/summary", h.GetDashboardSummary)

	return r, h, db
}

func generateToken(t *testing.T, secret string, userID uint) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	s, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

func executeRequest(r *gin.Engine, method, url string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func executeAuthenticatedRequest(r *gin.Engine, method, url, token string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

func createTestUser(t *testing.T, db *gorm.DB) uint {
	t.Helper()
	user := model.User{Username: "testuser", PasswordHash: "hash"}
	err := db.Create(&user).Error
	require.NoError(t, err)
	return user.ID
}
