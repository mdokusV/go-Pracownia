package controllers

import "github.com/gofiber/fiber/v2"

func SendWeb(c *fiber.Ctx) error {
	return c.Render("posts/index", fiber.Map{})
}
