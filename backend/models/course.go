package models

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Grade       string    `json:"grade"`
	Semester    string    `json:"semester"`
	Subject     string    `json:"subject"`
	Teacher     string    `json:"teacher"`
	Credits     int       `json:"credits"`
	ImageURL    string    `json:"imageURL"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Course) TableName() string {
	return "courses"
}
