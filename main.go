package main

import (
	"fmt"
	"log"

	config "go-booking-app/configs"
	_ "go-booking-app/docs"
	"go-booking-app/models"
	"go-booking-app/routes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Booking API
// @version 1.0
// @description This is a booking API server.
// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(&models.User{}, &models.Booking{}).Error; err != nil {
		log.Fatalf("migration error: %v", err)
	}

	r := routes.SetupRouter(db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
