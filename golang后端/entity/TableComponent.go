package entity

import "gorm.io/gorm"

type TableStyle struct {
	Width          string `json:"width"`
	MarginTop      string `json:"marginTop"`
	MarginBottom   string `json:"marginBottom"`
	BorderCollapse string `json:"borderCollapse"`
}

type TableComponent struct {
	gorm.Model
	ID              uint            `gorm:"primaryKey" json:"id"`
	Type            string          `json:"type"`
	Position        Position        `json:"position" gorm:"embedded"`
	TemplateId      uint            `json:"templateId"`
	Style           TableStyle      `json:"style" gorm:"embedded"`
	Title           string          `json:"title"`
	Layout          string          `json:"layout"`
	Rows            int             `json:"rows"`
	Columns         int             `json:"columns"`
	Data            StringMatrix    `json:"data" gorm:"type:json"`
	HeaderRow       bool            `json:"headerRow"`
	Borders         bool            `json:"borders"`
	TitleAlignment  string          `json:"titleAlignment"`
	TitleFontFamily string          `json:"titleFontFamily"`
	TitleFontSize   int             `json:"titleFontSize"`
	ColumnWidths    IntArray        `json:"columnWidths" gorm:"type:json"`
	Fields          TableFieldArray `json:"fields" gorm:"type:json"`
}

type TableField struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey" json:"id"`
	Label             string `json:"label"`
	Field             string `json:"field"`
	Type              string `json:"type"`
	Required          bool   `json:"required"`
	Value             string `json:"value"`
	Description       string `json:"description"`
	ContentComponents string `json:"contentComponents"`
}
