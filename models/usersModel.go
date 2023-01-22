package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"not null; size:128"`
	DateOfBirth string `gorm:"not null; size:128"`
	Surname     string `gorm:"not null; size:128"`
	Login       string `gorm:"not null; size:128; unique"`
	Password    string `gorm:"not null; size:128"`
	RoleID      uint   `gorm:"not null; foreignkey:RoleID"`
}
