package entity

import "gorm.io/gorm"

type CoverTableComponent struct {
	gorm.Model
	ID         uint                 `gorm:"primaryKey" json:"id"`
	Type       string               `json:"type"`
	Position   Position             `json:"position" gorm:"embedded"`
	TemplateId uint                 `json:"templateId"`
	Style      CoverTableStyle      `json:"style" gorm:"embedded"`
	Fields     CoverTableFieldArray `json:"fields" gorm:"type:json"`
	Settings   CoverTableSettings   `json:"settings" gorm:"embedded"`
}
type CoverTableField struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey" json:"id"`
	Label       string `json:"label"`
	Value       string `json:"value"`
	Width       int    `json:"width"`
	ComponentID uint   `json:"componentId"`
}
type CoverTableStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}

type CoverTableSettings struct {
	FontSize    int                   `json:"fontSize"`
	FontFamily  string                `json:"fontFamily"`
	LabelStyle  CoverTableLabelStyle  `json:"labelStyle" gorm:"embedded"`
	ValueStyle  CoverTableValueStyle  `json:"valueStyle" gorm:"embedded"`
	BorderStyle CoverTableBorderStyle `json:"borderStyle" gorm:"embedded"`
	CellStyle   CoverTableCellStyle   `json:"cellStyle" gorm:"embedded"`
	HeaderStyle CoverTableHeaderStyle `json:"headerStyle" gorm:"embedded"`
	LineHeight  float64               `json:"lineHeight"`
}

type CoverTableCellStyle struct {
	Alignment string `json:"alignment"`
	Padding   string `json:"padding"`
}

type CoverTableHeaderStyle struct {
	Bold      bool   `json:"bold"`
	Alignment string `json:"alignment"`
}

type CoverTableBorderStyle struct {
	Top              bool `json:"top"`
	Bottom           bool `json:"bottom"`
	Left             bool `json:"left"`
	Right            bool `json:"right"`
	InsideHorizontal bool `json:"insideHorizontal"`
	InsideVertical   bool `json:"insideVertical"`
}

type CoverTableLabelStyle struct {
	Bold      bool   `json:"bold"`
	Alignment string `json:"alignment"`
}

type CoverTableValueStyle struct {
	Alignment string `json:"alignment"`
}
