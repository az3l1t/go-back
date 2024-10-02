package controllers

import (
	"net/http"

	"go-booking-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UserRequestBooking struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CreateBooking создает новое бронирование
// @Summary Create a booking
// @Description Create a new booking. Use JWT token in headers.
// @Tags bookings
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param booking body UserRequestBooking true "Booking object"
// @Success 201 {object} models.Booking
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings [post]
func (bc *BookingController) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if booking.StartTime.IsZero() || booking.EndTime.IsZero() {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "StartTime and EndTime are required"})
		return
	}

	userID := c.MustGet("userID").(float64)
	booking.UserID = uint(userID)

	if err := bc.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	if err := bc.DB.Model(&booking).Related(&booking.User, "UserID").Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not load user data"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// DeleteBooking удаляет бронирование
// @Summary Delete a booking
// @Description Delete a booking by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "Booking ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings/{id} [delete]
func (bc *BookingController) DeleteBooking(c *gin.Context) {
	var booking models.Booking
	if err := bc.DB.Where("id = ?", c.Param("id")).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Booking not found"})
		return
	}

	if err := bc.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Error: "Booking deleted"})
}

// ListBookings получает список всех бронирований
// @Summary List bookings
// @Description Get all bookings
// @Tags bookings
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} models.Booking
// @Failure 500 {object} controllers.ErrorResponse
// @Router /bookings [get]
func (bc *BookingController) ListBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := bc.DB.Preload("User").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
