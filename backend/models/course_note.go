package models

import "time"

// CourseNote 课程笔记模型
type CourseNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null" json:"courseId"`
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`
	Content   string    `gorm:"type:text" json:"content"`
	AuthorID  uint      `gorm:"not null" json:"authorId"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	IsPublic  bool      `gorm:"default:false" json:"isPublic"` // 是否公开分享
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (CourseNote) TableName() string {
	return "course_notes"
}

