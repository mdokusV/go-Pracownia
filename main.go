package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/mdokusV/go-Pracownia/controllers"
	"github.com/mdokusV/go-Pracownia/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToMySQL()
}

func main() {
	// Load templates
	engine := html.New("./views", ".html")

	// Setup app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Config app
	app.Static("/", "./public")

	app.Post("/userAdd", controllers.UserCreate)
	app.Get("/ShowAllUsers", controllers.UserShowAll)

	app.Get("/", controllers.SendWeb)
	app.Listen(":" + os.Getenv("PORT")) // listen and serve on 0.0.0.0:3000
}
