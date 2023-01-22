package main

import (
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToMySQL()
}

func main() {
	// Migrate the schema

	//DropTable
	initializers.DB.Migrator().DropTable(&models.Role{}, &models.User{})
	initializers.DB.AutoMigrate(&models.User{}, &models.Role{})
	initializers.DB.Exec("ALTER TABLE users ADD FOREIGN KEY (role_id) REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE;")
}
