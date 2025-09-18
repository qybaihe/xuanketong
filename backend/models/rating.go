package models

import (
	"time"

	"gorm.io/gorm"
)

type Rating struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"userId"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	CourseID   uint      `json:"courseId"`
	Score      float64   `json:"score"`
	Difficulty float64   `json:"difficulty" gorm:"default:0"`
	Usefulness float64   `json:"usefulness" gorm:"default:0"`
	Teaching   float64   `json:"teaching" gorm:"default:0"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Rating) TableName() string {
	return "ratings"
}

// BeforeCreate 在创建记录前的钩子
func (r *Rating) BeforeCreate(tx *gorm.DB) error {
	// 验证总体评分
	if r.Score < 1 {
		r.Score = 1
	} else if r.Score > 5 {
		r.Score = 5
	}

	// 验证课程难度评分
	if r.Difficulty < 1 {
		r.Difficulty = 1
	} else if r.Difficulty > 5 {
		r.Difficulty = 5
	}

	// 验证实用性评分
	if r.Usefulness < 1 {
		r.Usefulness = 1
	} else if r.Usefulness > 5 {
		r.Usefulness = 5
	}

	// 验证教学质量评分
	if r.Teaching < 1 {
		r.Teaching = 1
	} else if r.Teaching > 5 {
		r.Teaching = 5
	}

	return nil
}

// BeforeUpdate 在更新记录前的钩子
func (r *Rating) BeforeUpdate(tx *gorm.DB) error {
	// 验证总体评分
	if r.Score < 1 {
		r.Score = 1
	} else if r.Score > 5 {
		r.Score = 5
	}

	// 验证课程难度评分
	if r.Difficulty < 1 {
		r.Difficulty = 1
	} else if r.Difficulty > 5 {
		r.Difficulty = 5
	}

	// 验证实用性评分
	if r.Usefulness < 1 {
		r.Usefulness = 1
	} else if r.Usefulness > 5 {
		r.Usefulness = 5
	}

	// 验证教学质量评分
	if r.Teaching < 1 {
		r.Teaching = 1
	} else if r.Teaching > 5 {
		r.Teaching = 5
	}

	return nil
}
