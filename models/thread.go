package models

import "time"

type Thread struct {
	BaseModel
	Text           string    `json:"text" binding:"required"`
	BumpedOn       time.Time `json:"bumped_on" gorm:"default:now()"`
	Reported       bool      `json:"reported" gorm:"default:false"`
	DeletePassword string    `json:"delete_password" gorm:"not null" binding:"required"`
	Replies        []Reply   `json:"replies,omitempty"`
}
