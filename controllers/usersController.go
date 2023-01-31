package controllers

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mdokusV/go-Pracownia/helpers"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const CLUSTER_SIZE = 10

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
	empty := initializers.DB.Where(&models.User{Login: body.Login}).First(&exists).Error
	if empty == nil {
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

	//Show for specific role type
	client := c.Locals("userCookie").(models.User)
	switch client.RoleID {
	case 1:
		//check if there are enough users to show
		firstUserToFind := (body.PageNumber - 1) * CLUSTER_SIZE
		var maxCount int64
		initializers.DB.Model(&models.User{}).Count(&maxCount)
		// initializers.DB.Model(&models.User{}).Count(&maxCount)
		if maxCount < int64(firstUserToFind) {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}

		// Find Users
		var users []models.User
		err := initializers.DB.Limit(CLUSTER_SIZE).Offset(firstUserToFind).Find(&users).Error
		if len(users) == 0 {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}
		if err != nil {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}

		// Prepare response JSON
		type sendJson struct {
			OrderNumber int
			Name        string
			Surname     string
			DateOfBirth string
			RoleName    string
		}
		var sendJsons []sendJson

		// Fill data into JSON
		startingOrderNumber := firstUserToFind + 1
		for _, user := range users {
			var role models.Role
			err := initializers.DB.Where(&models.Role{RoleID: user.RoleID}).Find(&role).Error //Find role with ID
			if err != nil {
				fmt.Println(err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
			}
			newUser := sendJson{
				OrderNumber: startingOrderNumber,
				Name:        user.Name,
				Surname:     user.Surname,
				DateOfBirth: user.DateOfBirth,
				RoleName:    role.Name,
			}
			startingOrderNumber++
			sendJsons = append(sendJsons, newUser)
		}

		//Send it
		return c.Status(fiber.StatusOK).JSON(sendJsons)

	case 2:
		//check if there are enough users to show
		firstUserToFind := (body.PageNumber - 1) * CLUSTER_SIZE
		var maxCount int64
		initializers.DB.Model(&models.User{}).Count(&maxCount)
		if maxCount < int64(firstUserToFind) {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}

		// Find Users
		var users []models.User
		err := initializers.DB.Limit(CLUSTER_SIZE).Offset(firstUserToFind).Find(&users).Error
		if len(users) == 0 {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}
		if err != nil {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}

		// Prepare response JSON
		type sendJson struct {
			OrderNumber int
			Name        string
			Surname     string
			DateOfBirth string
			RoleName    string
			Login       string
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}
		var sendJsons []sendJson

		// Fill data into JSON
		startingOrderNumber := firstUserToFind + 1
		for _, user := range users {
			var role models.Role
			err := initializers.DB.Where(&models.Role{RoleID: user.RoleID}).Find(&role).Error //Find role with ID
			if err != nil {
				fmt.Println(err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
			}
			newUser := sendJson{
				OrderNumber: startingOrderNumber,
				Name:        user.Name,
				Surname:     user.Surname,
				DateOfBirth: user.DateOfBirth,
				RoleName:    role.Name,
				Login:       user.Login,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
			}
			startingOrderNumber++
			sendJsons = append(sendJsons, newUser)
		}

		// Send it
		return c.Status(fiber.StatusOK).JSON(sendJsons)

	case 3:
		//check if there are enough users to show
		firstUserToFind := (body.PageNumber - 1) * CLUSTER_SIZE
		var maxCount int64
		initializers.DB.Model(&models.User{}).Unscoped().Count(&maxCount)
		if maxCount < int64(firstUserToFind) {
			return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}

		// Find Users
		var users []models.User
		initializers.DB.Limit(CLUSTER_SIZE).Unscoped().Offset(firstUserToFind).Find(&users)

		// Prepare response JSON
		type sendJson struct {
			OrderNumber int
			Name        string
			Surname     string
			DateOfBirth string
			RoleName    string
			Login       string
			Password    string
			gorm.Model
		}
		var sendJsons []sendJson

		// Fill data into JSON
		startingOrderNumber := firstUserToFind + 1
		for _, user := range users {
			var role models.Role
			err := initializers.DB.Where(&models.Role{RoleID: user.RoleID}).Find(&role).Error //Find role with ID
			if err != nil {
				fmt.Println(err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
			}
			newUser := sendJson{
				OrderNumber: startingOrderNumber,
				Name:        user.Name,
				Surname:     user.Surname,
				DateOfBirth: user.DateOfBirth,
				RoleName:    role.Name,
				Login:       user.Login,
				Password:    user.Password,
				Model:       user.Model,
			}
			startingOrderNumber++
			sendJsons = append(sendJsons, newUser)
		}

		// Send It
		return c.Status(fiber.StatusOK).JSON(sendJsons)
	default:
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}

func UserChangeRole(c *fiber.Ctx) error {
	// Get data from request
	var body struct {
		Login    string `validate:"required,email,min=6,max=32"`
		RoleName string `validate:"required,min=3,max=32"`
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

	// Validate User Admin
	client := c.Locals("userCookie").(models.User)
	if client.RoleID != 3 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var err error
	// Find role
	var role models.Role
	err = initializers.DB.Where(&models.Role{Name: body.RoleName}).First(&role).Error
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	// Find user
	var user models.User
	err = initializers.DB.Where(&models.User{Login: body.Login}).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	// Change it
	user.RoleID = role.RoleID
	err = initializers.DB.Save(&user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func MaxPages(c *fiber.Ctx) error {
	// Show for specific role type
	client := c.Locals("userCookie").(models.User)

	var maxCount int64
	var err error

	switch client.RoleID {
	case 1:
		err = initializers.DB.Model(&models.User{}).Count(&maxCount).Error
	case 2:
		err = initializers.DB.Model(&models.User{}).Count(&maxCount).Error
	case 3:
		err = initializers.DB.Model(&models.User{}).Unscoped().Count(&maxCount).Error
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}
	maxPages := math.Ceil(float64(maxCount) / CLUSTER_SIZE)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"MaxPages": maxPages})
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
	fmt.Println("XDDD")
	return c.Redirect("/login", fiber.StatusMovedPermanently)
}

func UserDelete(c *fiber.Ctx) error {
	// Get data from request
	var body struct {
		Login string `validate:"required,email,min=6,max=32"`
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

	// Check if Admin
	client := c.Locals("userCookie").(models.User)

	if client.RoleID != 3 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// Check if user exists
	exists := initializers.DB.Where(&models.User{Login: body.Login}).First(&models.User{}).Error
	if exists != nil {
		return c.SendStatus(fiber.StatusExpectationFailed)
	}

	// Delete the user
	result := initializers.DB.Where(&models.User{Login: body.Login}).Delete(&models.User{}).Error
	if result != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(result)
	}
	return c.SendStatus(fiber.StatusOK)
}
