package models

import "time"

type Thread struct {
	BaseModel
	Text           string    `json:"text"`
	BumpedOn       time.Time `json:"bumped_on" gorm:"default:now()"`
	Reported       bool      `json:"reported" gorm:"default:false"`
	DeletePassword string    `json:"delete_password" gorm:"not null"`
	Replies        []Reply   `json:"replies,omitempty"`
}
