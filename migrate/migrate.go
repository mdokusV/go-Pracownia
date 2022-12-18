package main

import (
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToSQLite()
}

func main() {
	// Migrate the schema
	initializers.DB.AutoMigrate(&models.Post{})
}
