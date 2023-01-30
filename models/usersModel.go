package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"not null; size:64"`
	DateOfBirth string `gorm:"not null; type: date"`
	Surname     string `gorm:"not null; size:64"`
	Login       string `gorm:"not null; size:64; unique"`
	Password    string `gorm:"not null; size:128"`
	RoleID      uint   `gorm:"not null; foreignkey:RoleID"`
}
