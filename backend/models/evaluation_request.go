package models

import "time"

// EvaluationRequest represents a user's request for course evaluation.
type EvaluationRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CourseID  uint      `gorm:"not null" json:"courseId"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // e.g., 'pending', 'closed'
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (EvaluationRequest) TableName() string {
	return "evaluation_requests"
}
