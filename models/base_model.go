package models

import "time"

type BaseModel struct {
	ID      uint      `json:"id" gorm:"primary_key"`
	Created time.Time `json:"created" gorm:"default:now()"`
}
