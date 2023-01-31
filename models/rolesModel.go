package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleID uint   `gorm:"not null; size:128; unique"`
	Name   string `gorm:"not null; size:128; unique"`
}
