package entity

import "gorm.io/gorm"

type Template struct {
	gorm.Model
	ID               uint             `gorm:"primaryKey" json:"id,omitempty"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	UserId           uint             `json:"userId"`
	DocumentSettings DocumentSettings `json:"documentSettings" gorm:"embedded"`
}

type DocumentSettings struct {
	PageSize PageSize `json:"pageSize" gorm:"embedded"`
	Margins  Margins  `json:"margins" gorm:"embedded"`
}

type PageSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Margins struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

type TitleStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
	TextAlign    string `json:"textAlign"`
	FontWeight   string `json:"fontWeight"`
}

type TitleComponent struct {
	gorm.Model
	ID         uint       `gorm:"primaryKey" json:"id"`
	Type       string     `json:"type"`
	Position   Position   `json:"position" gorm:"embedded"`
	TemplateId uint       `json:"templateId"`
	Style      TitleStyle `json:"style" gorm:"embedded"`
	Content    string     `json:"content"`
	FontSize   int        `json:"fontSize"`
	FontFamily string     `json:"fontFamily"`
	Alignment  string     `json:"alignment"`
	Bold       bool       `json:"bold"`
}

type TextStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}
type TextComponent struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `json:"type"`
	Position   Position  `json:"position" gorm:"embedded"`
	TemplateId uint      `json:"templateId"`
	Style      TextStyle `json:"style" gorm:"embedded"`
	Content    string    `json:"content"`
	FontSize   int       `json:"fontSize"`
	FontFamily string    `json:"fontFamily"`
	Alignment  string    `json:"alignment"`
	LineHeight int       `json:"lineHeight"`
}

type ImageStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}

type ImageComponent struct {
	gorm.Model
	ID                uint       `gorm:"primaryKey" json:"id"`
	Type              string     `json:"type"`
	Position          Position   `json:"position" gorm:"embedded"`
	TemplateId        uint       `json:"templateId"`
	Style             ImageStyle `json:"style" gorm:"embedded"`
	ImageUrl          string     `json:"imageUrl"`
	ImagePreviewError bool       `json:"imagePreviewError"`
	ImageAlignment    string     `json:"imageAlignment"`
	MaxWidth          int        `json:"maxWidth"`
	MaxHeight         string     `json:"maxHeight"`
	ShowCaption       bool       `json:"showCaption"`
	Caption           string     `json:"caption"`
	CaptionAlignment  string     `json:"captionAlignment"`
	CaptionFontSize   int        `json:"captionFontSize"`
	CaptionFontFamily string     `json:"captionFontFamily"`
	CaptionLineHeight float64    `json:"captionLineHeight"`
}
type TocStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}
type TocComponent struct {
	gorm.Model
	ID              uint         `gorm:"primaryKey" json:"id"`
	Type            string       `json:"type"`
	Position        Position     `json:"position" gorm:"embedded"`
	TemplateId      uint         `json:"templateId"`
	Style           TocStyle     `json:"style" gorm:"embedded"`
	Title           string       `json:"title"`
	TitleFontFamily string       `json:"titleFontFamily"`
	TitleFontSize   int          `json:"titleFontSize"`
	FontSize        int          `json:"fontSize"`
	FontFamily      string       `json:"fontFamily"`
	LineHeight      float64      `json:"lineHeight"`
	Alignment       string       `json:"alignment"`
	Items           TocItemArray `json:"items" gorm:"type:json"`
}
type TocItem struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey" json:"id"`
	Number string `json:"number"`
	Title  string `json:"title"`
	Page   string `json:"page"`
}

type ListComponent struct {
	gorm.Model
	ID         uint          `gorm:"primaryKey" json:"id"`
	Type       string        `json:"type"`
	Position   Position      `json:"position" gorm:"embedded"`
	TemplateId uint          `json:"templateId"`
	Style      ListStyle     `json:"style" gorm:"embedded"`
	Items      ListItemArray `json:"items" gorm:"type:json"`
	ListType   string        `json:"listType"`
	FontSize   int           `json:"fontSize"`
	FontFamily string        `json:"fontFamily"`
	LineHeight float64       `json:"lineHeight"`
	Alignment  string        `json:"alignment"`
}
type ListItem struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"id"`
	Content string `json:"content"`
}
type ListStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ReferenceStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}
type ReferenceComponent struct {
	gorm.Model
	ID              uint               `gorm:"primaryKey" json:"id"`
	Type            string             `json:"type"`
	Position        Position           `json:"position" gorm:"embedded"`
	TemplateId      uint               `json:"templateId"`
	Style           ReferenceStyle     `json:"style" gorm:"embedded"`
	Title           string             `json:"title"`
	TitleFontFamily string             `json:"titleFontFamily"`
	TitleFontSize   int                `json:"titleFontSize"`
	FontSize        int                `json:"fontSize"`
	FontFamily      string             `json:"fontFamily"`
	LineHeight      float64            `json:"lineHeight"`
	Alignment       string             `json:"alignment"`
	Items           ReferenceItemArray `json:"items" gorm:"type:json"`
}

type ReferenceItem struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"id"`
	Content string `json:"content"`
}
type ScheduleStyle struct {
	Width        string `json:"width"`
	MarginTop    string `json:"marginTop"`
	MarginBottom string `json:"marginBottom"`
}
type ScheduleItem struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	Task      string `json:"task"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
type ScheduleComponent struct {
	gorm.Model
	ID         uint              `gorm:"primaryKey" json:"id"`
	Type       string            `json:"type"`
	Position   Position          `json:"position" gorm:"embedded"`
	TemplateId uint              `json:"templateId"`
	Style      ScheduleStyle     `json:"style" gorm:"embedded"`
	Items      ScheduleItemArray `json:"items" gorm:"type:json"`
	Settings   ScheduleSettings  `json:"settings" gorm:"embedded"`
}
type ScheduleSettings struct {
	FontSize    int         `json:"fontSize"`
	FontFamily  string      `json:"fontFamily"`
	LineHeight  float64     `json:"lineHeight"`
	BorderStyle BorderStyle `json:"borderStyle" gorm:"embedded"`
	CellStyle   CellStyle   `json:"cellStyle" gorm:"embedded"`
	HeaderStyle HeaderStyle `json:"headerStyle" gorm:"embedded"`
	LabelStyle  LabelStyle  `json:"labelStyle" gorm:"embedded"`
	ValueStyle  ValueStyle  `json:"valueStyle" gorm:"embedded"`
}
type BorderStyle struct {
	Top              bool `json:"top"`
	Bottom           bool `json:"bottom"`
	Left             bool `json:"left"`
	Right            bool `json:"right"`
	InsideHorizontal bool `json:"insideHorizontal"`
	InsideVertical   bool `json:"insideVertical"`
}
type CellStyle struct {
	Alignment string `json:"alignment"`
	Padding   string `json:"padding"`
}
type HeaderStyle struct {
	Bold      bool   `json:"bold"`
	Alignment string `json:"alignment"`
}
type LabelStyle struct {
	Bold      bool   `json:"bold"`
	Alignment string `json:"alignment"`
}
type ValueStyle struct {
	Alignment string `json:"alignment"`
}
