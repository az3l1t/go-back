package tests

import (
	"bytes"
	"encoding/json"
	"go-booking-app/controllers"
	"go-booking-app/models"
	"go-booking-app/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres password=Chemege1. dbname=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Booking{})

	router = routes.SetupRouter(db)

	m.Run()

	CleanupDB()
}

// Функция для очистки данных после каждого теста
func CleanupDB() {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM bookings")
}

func TestRegisterUser(t *testing.T) {
	user := controllers.UserRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response controllers.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "testuser", response.Username)
}

func TestLoginUser(t *testing.T) {
	user := controllers.UserRequest{
		Username: "testlogin",
		Password: "testpassword",
	}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tokenResponse controllers.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &tokenResponse)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenResponse.Token)
}
