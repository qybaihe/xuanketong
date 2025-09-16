package models

import "time"

type Rating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"userId"`
	CourseID  uint      `json:"courseId"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Rating) TableName() string {
	return "ratings"
}
