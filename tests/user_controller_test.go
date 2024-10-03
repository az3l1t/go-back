package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-booking-app/models"
	"go-booking-app/routes"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(db *gorm.DB) *gin.Engine {
	return routes.SetupRouter(db)
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=localhost port=5433 user=postgres dbname=test_db password=Chemege1. sslmode=disable")
	if err != nil {
		return nil, err
	}

	db.Exec("DELETE FROM users")
	db.AutoMigrate(&models.User{}, &models.Booking{})
	return db, nil
}

func TestRegisterUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	w := httptest.NewRecorder()
	reqBody := `{"username": "testuser1", "password": "password"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "testuser1", response["username"])
}

func TestLoginUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	user := models.User{
		Username: "testuser1",
		Password: "password",
	}
	db.Create(&user)

	w := httptest.NewRecorder()
	reqBody := `{"username": "testuser1", "password": "password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}
