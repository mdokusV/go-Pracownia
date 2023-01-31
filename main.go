package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/mdokusV/go-Pracownia/controllers"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/middleware"
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
	app.Use(logger.New())

	// Config app
	app.Static("/", "./public")

	// app.Get("/MainPage", controllers.MainPageSend)

	app.Get("/login", controllers.SendLoginPage)
	app.Post("/login", controllers.UserLogin)
	app.Post("/logout", controllers.UserLogout)

	app.Get("/register", controllers.SendRegisterPage)
	app.Post("/register", controllers.UserCreate)

	app.Get("/UserShowAll", middleware.CheckAuthCookie, controllers.UserShowAll)
	app.Post("/UserShow", middleware.CheckAuthCookie, controllers.UserShow)
	app.Get("/maxPages", middleware.CheckAuthCookie, controllers.MaxPages)
	app.Post("/changeUserRole", middleware.CheckAuthCookie, controllers.UserChangeRole)

	app.Delete("/UserDelete", middleware.CheckAuthCookie, controllers.UserDelete)

	app.Get("/", middleware.CheckAuthCookie, controllers.SendMainWeb)

	app.Listen(":" + os.Getenv("PORT")) // listen and serve on 0.0.0.0:3000
}
