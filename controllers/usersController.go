package controllers

import (
	"fmt"
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
)

func UserCreate(c *fiber.Ctx) error {
	// Get data from request
	var body struct {
		Name        string
		Surname     string
		DateOfBirth string
		Login       string
		RoleID      uint
		Password    string
	}
	c.BodyParser(&body)

	//Check if login is email
	_, err := mail.ParseAddress(body.Login)
	if err != nil {
		c.Status(400)
		fmt.Println(err)
		return err
	}

	//Create POST
	User := models.User{
		Name:        body.Name,
		Surname:     body.Surname,
		DateOfBirth: body.DateOfBirth,
		Login:       body.Login,
		RoleID:      body.RoleID,
		Password:    body.Password,
	}

	result := initializers.DB.Create(&User)
	if result.Error != nil {
		c.Status(400)
		fmt.Println(result.Error)
		return result.Error
	}

	//Return it

	return c.Status(200).JSON(User)

}

func UserShowAll(c *fiber.Ctx) error {
	//Get the posts
	var users []models.User
	initializers.DB.Find(&users)

	//Respond with them
	return c.Status(200).JSON(users)
}
