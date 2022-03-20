package services

import (
	"errors"

	"github.com/raulandre/anonboard/database"
	"github.com/raulandre/anonboard/models"
	"github.com/raulandre/anonboard/utils"
	"gorm.io/gorm"
)

type ThreadService interface {
	List(page int, pageSize int) (*[]models.Thread, error)
	GetById(id string) (*models.Thread, error)
	Create(t models.Thread) (*models.Thread, error)
	Report(id string) error
	DeleteWithPassword(id, password string) error
}

type threadService struct {
	db *gorm.DB
}

func NewThreadService(conn database.DatabaseConnection) ThreadService {
	return &threadService{db: conn.Get()}
}

func (ts *threadService) List(page int, pageSize int) (*[]models.Thread, error) {
	var t []models.Thread

	offset := 0
	if page > 0 {
		offset = page - 1
	}

	res := ts.db.
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Order("bumped_on DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&t)

	return &t, res.Error
}

func (ts *threadService) GetById(id string) (*models.Thread, error) {
	var t models.Thread
	res := ts.db.
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Where("id = ?", id).
		First(&t)
	return &t, res.Error
}

func (ts *threadService) Create(t models.Thread) (*models.Thread, error) {
	password, err := utils.HashPassword(t.DeletePassword)

	if err != nil {
		return nil, err
	}

	t.DeletePassword = password

	res := ts.db.Create(&t)
	return &t, res.Error
}

func (ts *threadService) Report(id string) error {
	return ts.db.Transaction(func(tx *gorm.DB) error {
		var t models.Thread
		if result := tx.Where("id = ?", id).First(&t); result.Error != nil {
			return result.Error
		}
		if result := tx.Model(&t).Where("id = ?", id).Update("reported", true); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (ts *threadService) DeleteWithPassword(id, password string) error {
	return ts.db.Transaction(func(tx *gorm.DB) error {
		var t models.Thread
		if result := tx.Where("id = ?", id).First(&t); result.Error != nil {
			return result.Error
		}
		if !utils.CheckPassword(password, t.DeletePassword) {
			return errors.New("incorrect password")
		}
		if result := tx.Model(&t).Where("id = ?", id).Update("reported", true); result.Error != nil {
			return result.Error
		}
		return nil
	})
}
