package routes

import (
	"go-booking-app/controllers"
	middleware "go-booking-app/filters"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	userController := controllers.UserController{DB: db}
	bookingController := controllers.BookingController{DB: db}

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	r.DELETE("/users/:id", userController.DeleteUser)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/bookings", bookingController.CreateBooking)
		auth.DELETE("/bookings/:id", bookingController.DeleteBooking)
		auth.GET("/bookings", bookingController.ListBookings)
	}

	return r
}
