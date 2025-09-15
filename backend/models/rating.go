package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	UserID   uint
	CourseID uint
	Score    float64
}
