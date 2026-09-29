package entity

import "gorm.io/gorm"

type Virtual struct {
	gorm.Model
	ID           uint         `json:"id"`
	TemplateId   uint         `json:"templateId" gorm:"column:template_id"`
	StudentId    uint         `json:"studentId" gorm:"column:student_id"`
	Difficulty   string       `json:"difficulty"`
	User         User         `json:"user" gorm:"foreignKey:StudentId;references:ID"`
	Seconds      int          `json:"seconds"`
	Score        int          `json:"score"`
	Description  string       `json:"description"`
	FakeTemplate FakeTemplate `json:"fake_template" gorm:"foreignKey:TemplateId;references:ID"`
}

type VirtualQuestion struct {
	gorm.Model
	ID                uint    `json:"id"`
	VirtualId         uint    `json:"virtualId" gorm:"column:virtual_id"`
	Question          string  `json:"question"`
	Answer            string  `json:"answer"`
	QuestionAudioPath string  `json:"question_audio_path"`
	AnswerAudioPath   string  `json:"answer_audio_path"`
	Virtual           Virtual `json:"virtual" gorm:"foreignKey:VirtualId;references:ID"`
}
