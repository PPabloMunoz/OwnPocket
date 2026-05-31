package service

import (
	"testing"

	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser_Success(t *testing.T) {
	svc, db := setupService(t)

	email := "test@example.com"
	err := svc.RegisterUser("testuser", "password123", &email)
	require.NoError(t, err)

	var user model.User
	err = db.Where("username = ?", "testuser").First(&user).Error
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", *user.Email)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
	assert.NoError(t, err)
}

func TestRegisterUser_NoEmail(t *testing.T) {
	svc, db := setupService(t)

	err := svc.RegisterUser("testuser", "password123", nil)
	require.NoError(t, err)

	var user model.User
	err = db.Where("username = ?", "testuser").First(&user).Error
	require.NoError(t, err)
	assert.Nil(t, user.Email)
}

func TestRegisterUser_EmptyEmail(t *testing.T) {
	svc, db := setupService(t)

	email := ""
	err := svc.RegisterUser("testuser", "password123", &email)
	require.NoError(t, err)

	var user model.User
	err = db.Where("username = ?", "testuser").First(&user).Error
	require.NoError(t, err)
	assert.Nil(t, user.Email)
}

func TestRegisterUser_DuplicateUsername(t *testing.T) {
	svc, _ := setupService(t)

	email := "test@example.com"
	err := svc.RegisterUser("testuser", "password123", &email)
	require.NoError(t, err)

	otherEmail := "other@example.com"
	err = svc.RegisterUser("testuser", "otherpass", &otherEmail)
	assert.Error(t, err)
}

func TestLogin_Success(t *testing.T) {
	svc, _ := setupService(t)

	email := "test@example.com"
	err := svc.RegisterUser("testuser", "password123", &email)
	require.NoError(t, err)

	token, err := svc.Login("testuser", "password123", "test-secret")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_InvalidPassword(t *testing.T) {
	svc, _ := setupService(t)

	email := "test@example.com"
	err := svc.RegisterUser("testuser", "password123", &email)
	require.NoError(t, err)

	token, err := svc.Login("testuser", "wrongpass", "test-secret")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _ := setupService(t)

	token, err := svc.Login("nonexistent", "password123", "test-secret")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
}
