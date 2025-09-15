package controllers

import (
	"fmt"
	"time"
	"os"
	"project/models"
	"project/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	result :=  database.DB.Find(&users)     //fetching all users from database and map them in to users//
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	return c.JSON(users) //successfully return in json//

}
func RegisterUser(c *fiber.Ctx) error {
	var input models.User                     //like slice to user info//
	if err:=c.BodyParser(&input);err!=nil{       //read request and store in var//
		return c.Status(422).JSON(fiber.Map{"error":"Invalid JSON Input"})
}
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(422).JSON(fiber.Map{"error": "Missing required field"}) //422 iput valid but missing some values//
	}
	var exists models.User
	if err := database.DB.Where("email = ?", input.Email).First(&exists).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already registered"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)  //to cover password//
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     "user", // default role
	}
	if err := database.DB.Create(&user).Error; err != nil {              //store user in db//
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}


func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		fmt.Println(input)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})  //checking db pass and input pass same or not//
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"role":     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}


