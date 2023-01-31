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
	err := initializers.DB.Where(&models.Role{RoleID: client.RoleID}).First(&role).Error
	if err != nil {
		return c.Render("/login", fiber.Map{})
	}

	return c.Render("posts/index"+cases.Title(language.Und).String(role.Name)+"MainPage", fiber.Map{})
}

func SendLoginPage(c *fiber.Ctx) error {
	c.ClearCookie("Authorization")
	return c.Render("posts/indexLoginPage", fiber.Map{})
}

func SendRegisterPage(c *fiber.Ctx) error {
	c.ClearCookie("Authorization")
	return c.Render("posts/indexRegistrationPage", fiber.Map{})
}
