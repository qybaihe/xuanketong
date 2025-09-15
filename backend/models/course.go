package models

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Name        string
	Description string
	Grade       string
	Semester    string
	Subject     string
	Teacher     string
	Credits     int
	ImageURL    string
}
