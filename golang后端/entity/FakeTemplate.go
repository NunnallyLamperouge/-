package entity

import (
	"gorm.io/gorm"
)

type FakeTemplate struct {
	gorm.Model
	ID          uint   `json:"id" `
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	UserId      uint   `json:"userId"`
	User        User   `json:"user" gorm:"foreignKey:UserId"`
	Content     string `json:"content" gorm:"type:longtext"`
	Code        string `json:"code" gorm:"type:varchar(30)"`
	Slug        string `json:"slug" gorm:"type:varchar(30)"`
}

type TemplateHistory struct {
	gorm.Model
	ID         uint   `json:"id"`
	TemplateId uint   `json:"templateId"`
	Content    string `json:"content" gorm:"type:longtext"`
}
