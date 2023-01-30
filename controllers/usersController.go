package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mdokusV/go-Pracownia/helpers"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(c *fiber.Ctx) error {
	// Get data from request
	var body struct {
		Name        string `validate:"required,min=3,max=32"`
		Surname     string `validate:"required,min=3,max=32"`
		DateOfBirth string `validate:"required,dateformat"`
		Login       string `validate:"required,email,min=6,max=32"`
		Password    string `validate:"required,min=6,max=32"`
	}

	bindBody := c.BodyParser(&body)
	if bindBody != nil {
		fmt.Println(bindBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Failed to read body"})
	}

	// Validate struct
	errors := helpers.ValidateStruct(&body)
	if errors != nil {
		fmt.Println(errors)
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	//check if login already exists
	var exists models.User
	// initializers.DB.Where(&models.User{Login: body.Login}).First(&exists)
	initializers.DB.First(&exists, body.Login)
	if exists.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Already exists"})
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Failed to hash password"})
	}

	//Create User
	User := models.User{
		Name:        body.Name,
		Surname:     body.Surname,
		DateOfBirth: body.DateOfBirth,
		Login:       body.Login,
		RoleID:      1,
		Password:    string(hash),
	}

	result := initializers.DB.Create(&User)
	if result.Error != nil {
		fmt.Println(result.Error)
		return c.Status(fiber.StatusConflict).JSON(result.Error)
	}

	//Return it
	return c.Redirect("/login", fiber.StatusMovedPermanently)
}

func UserShowAll(c *fiber.Ctx) error {
	//Get the posts
	var users []models.User
	initializers.DB.Find(&users)

	//Respond with them
	return c.Status(200).JSON(users)
}

func UserShow(c *fiber.Ctx) error {
	const CLUSTER_SIZE = 10
	// Get data from request
	var body struct {
		PageNumber int `validate:"required,max=256"`
	}

	bindBody := c.BodyParser(&body)
	if bindBody != nil {
		fmt.Println(bindBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Failed to read body"})
	}

	// Validate struct
	errors := helpers.ValidateStruct(&body)
	if errors != nil {
		fmt.Println(errors)
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	//check if there are enough users to show
	firstUserToFind := (body.PageNumber - 1) * CLUSTER_SIZE
	var maxCount int64
	initializers.DB.Model(&models.User{}).Count(&maxCount)
	if maxCount < int64(firstUserToFind) {
		return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
	}

	// Find Users
	var users []models.User
	initializers.DB.Limit(CLUSTER_SIZE).Offset(firstUserToFind).Find(&users)

	//Show for specific role type
	client := c.Locals("userCookie").(models.User)

	if client.RoleID == 1 {
		type sendJson struct {
			ID          uint
			Name        string
			DateOfBirth string
			Surname     string
			RoleID      uint
		}
		var sendJsons []sendJson

		for _, user := range users {
			newUser := sendJson{
				ID:          user.ID,
				Name:        user.Name,
				Surname:     user.DateOfBirth,
				DateOfBirth: user.DateOfBirth,
				RoleID:      user.RoleID,
			}
			sendJsons = append(sendJsons, newUser)
		}
		return c.Status(fiber.StatusOK).JSON(sendJsons)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func UserLogin(c *fiber.Ctx) error {
	// Get the email and pass off request
	var body struct {
		Login    string `validate:"required,email,min=6,max=32"`
		Password string `validate:"required,min=6,max=32"`
	}

	bindBody := c.BodyParser(&body)
	if bindBody != nil {
		fmt.Println(bindBody)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Failed to read body"})
	}

	// Validate struct
	errors := helpers.ValidateStruct(&body)
	if errors != nil {
		fmt.Println(errors)
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Look up requested user
	var user models.User
	initializers.DB.Where(&models.User{Login: body.Login}).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid email"})
	}

	// Compare pass with saved pass-hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Invalid password"})
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Failed to create token"})
	}

	// Clear cookie
	c.ClearCookie("Authorization")

	// Send IT
	cookie := new(fiber.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	cookie.HTTPOnly = true
	cookie.Secure = false
	c.Cookie(cookie)
	return c.Redirect("/", fiber.StatusMovedPermanently)
}

func UserLogout(c *fiber.Ctx) error {
	c.ClearCookie("Authorization")
	return c.Redirect("/login", fiber.StatusMovedPermanently)
}
