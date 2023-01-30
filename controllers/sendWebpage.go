package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SendMainWeb(c *fiber.Ctx) error {
	// Check User Role
	client := c.Locals("userCookie").(models.User)

	var role models.Role
	err := initializers.DB.First(&role, client.RoleID).Error
	if err != nil {
		return c.Render("/login", fiber.Map{})
	}
	//TODO main web page
	return c.Render("posts/index"+cases.Title(language.Und).String(role.Name)+"MainPage", fiber.Map{})
}

//TODO login web page

func SendLoginPage(c *fiber.Ctx) error {
	return c.Render("posts/indexLoginPage", fiber.Map{})
}

//TODO registration web page

func SendRegisterPage(c *fiber.Ctx) error {
	c.ClearCookie("Authorization")
	return c.Render("posts/indexRegistrationPage", fiber.Map{})
}
