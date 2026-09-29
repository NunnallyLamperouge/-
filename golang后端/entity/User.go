package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey" json:"id"`
	Username    string `json:"username" gorm:"type:varchar(50)"`
	Password    string `json:"password" gorm:"type:varchar(1500)"`
	Email       string `json:"email" gorm:"type:varchar(50)"`
	Avatar      string `json:"avatar" gorm:"type:varchar(80)"`
	Tokens      int    `json:"tokens" gorm:"type:MEDIUMINT"`
	Permission  uint8  `json:"permission" gorm:"type:tinyint"`
	Type        uint8  `json:"type" gorm:"type:tinyint"`
	Group       uint8  `json:"group" gorm:"type:tinyint"`
	Apartment   string `json:"apartment" gorm:"type:varchar(50)"`
	Profession  string `json:"profession" gorm:"type:varchar(50)"`
	WorkspaceId string `json:"workspaceId" gorm:"type:varchar(50)"`
}
