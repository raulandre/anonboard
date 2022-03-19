package models

type Reply struct {
	BaseModel
	Text           string `json:"text"`
	Reported       bool   `json:"reported" gorm:"default:false"`
	DeletePassword string `json:"delete_password" gorm:"not null"`
	ThreadID       uint   `json:"thread_id"`
}
