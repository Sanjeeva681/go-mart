package database

import (
	"fmt"
	"log"
	"os"

	"project/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB
func ConnectDb(){
	dsn:=os.Getenv(("DB_DSN"))
	db,err:=gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	}),&gorm.Config{})
	if err!=nil{
		log.Fatal("failed to connect to database",err)
	}
	DB=db
	DB.AutoMigrate(&models.User{},&models.Product{},&models.Cart{},&models.CartItem{},&models.Order{},&models.OrderItem{})
	fmt.Println("Database connected succesfully")
}