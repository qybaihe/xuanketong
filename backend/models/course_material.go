package models

import "time"

// CourseMaterial 课程资料模型
type CourseMaterial struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CourseID       uint      `gorm:"not null" json:"courseId"`
	Course         Course    `gorm:"foreignKey:CourseID" json:"course"`
	MaterialName   string    `gorm:"not null" json:"materialName"`
	Type           string    `gorm:"not null" json:"type"`    // cloud, note
	Service        string    `gorm:"not null" json:"service"` // 云服务提供商
	Link           string    `gorm:"not null" json:"link"`    // 资料链接
	AccessPassword string    `json:"accessPassword"`          // 访问密码
	Description    string    `json:"description"`             // 资料描述
	UploaderID     uint      `gorm:"not null" json:"uploaderId"`
	Uploader       User      `gorm:"foreignKey:UploaderID" json:"uploader"`
	UploadTime     time.Time `json:"uploadTime"`
	DownloadCount  int       `gorm:"default:0" json:"downloadCount"`
	IsActive       bool      `gorm:"default:true" json:"isActive"` // 软删除标记
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (CourseMaterial) TableName() string {
	return "course_materials"
}
