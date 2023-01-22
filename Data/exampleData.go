package main

import (
	//"fmt"
	//"math/rand"

	"fmt"
	"math/rand"
	"time"

	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
	"gorm.io/gorm"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToMySQL()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	createRole("user", 1)
	createRole("moderator", 2)
	createRole("admin", 3)

	names := []string{"Emma", "Liam", "Olivia", "Noah", "Ava", "Ethan", "Isabella", "Mason", "Sophia", "Jacob"}
	surnames := []string{"Smith", "Johnson", "Williams", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Perez"}
	domains := []string{"gmail.com", "yahoo.com", "wp.com"}

	for i := 0; i < 20; i++ {
		createUser(
			names[rand.Intn(len(names))],
			surnames[rand.Intn(len(names))],
			generateDateOfBirth(),
			generateEmail(domains),
			generatePassword(10),
			rand.Intn(3)+1,
		)
	}
}

func generateEmail(domains []string) string {
	rand.Seed(time.Now().UnixNano())

	var username string
	for i := 0; i < 8; i++ {
		username += string(rand.Intn(26) + 'a')
	}

	domainIndex := rand.Intn(len(domains))
	domain := domains[domainIndex]

	email := username + "@" + domain

	return email
}

func generateDateOfBirth() string {
	rand.Seed(time.Now().UnixNano())

	day := rand.Intn(31) + 1
	month := rand.Intn(12) + 1
	year := rand.Intn(120) + 1901

	date := fmt.Sprintf("%02d/%02d/%04d", day, month, year)

	return date
}

func generatePassword(passwordLength int) string {
	const (
		passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	)
	password := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		password[i] = passwordChars[rand.Intn(len(passwordChars))]
	}
	return string(password)
}

func createUser(name, surname, dateOfBirth, login, password string, roleID int) {
	var User = models.User{
		Name:        name,
		Surname:     surname,
		DateOfBirth: dateOfBirth,
		Login:       login,
		RoleID:      uint(roleID),
		Password:    password,
	}

	var result = initializers.DB.Create(&User)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func createRole(name string, number uint) {
	var Role = models.Role{
		Name: name,
		Model: gorm.Model{
			ID: number,
		},
	}

	var result = initializers.DB.Create(&Role)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}
