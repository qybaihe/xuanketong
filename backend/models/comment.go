package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CourseID  uint      `json:"courseId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Comment) TableName() string {
	return "comments"
}
