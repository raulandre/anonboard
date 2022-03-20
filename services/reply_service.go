package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/raulandre/anonboard/database"
	"github.com/raulandre/anonboard/models"
	"github.com/raulandre/anonboard/utils"
	"gorm.io/gorm"
)

type ReplyService interface {
	ListByThreadId(tid, page, pageSize int) (*[]models.Reply, error)
	GetById(id string) (*models.Reply, error)
	Create(tid string, r models.Reply) (*models.Reply, error)
	Report(id string) error
	DeleteWithPassword(id, password string) error
}

type replyService struct {
	db *gorm.DB
}

func NewReplyService(conn database.DatabaseConnection) ReplyService {
	return &replyService{db: conn.Get()}
}

func (rs *replyService) ListByThreadId(tid, page, pageSize int) (*[]models.Reply, error) {
	var r []models.Reply
	offset := 0
	if page > 0 {
		offset = page - 1
	}
	res := rs.db.Where("thread_id = ?", tid).Limit(pageSize).Offset(offset).Find(&r)
	return &r, res.Error
}

func (rs *replyService) GetById(id string) (*models.Reply, error) {
	var r models.Reply
	res := rs.db.
		Where("id = ?", id).
		First(&r)
	return &r, res.Error
}

func (rs *replyService) Create(tid string, r models.Reply) (*models.Reply, error) {
	password, err := utils.HashPassword(r.DeletePassword)
	if err != nil {
		return nil, err
	}

	threadId, err := strconv.Atoi(tid)
	if err != nil {
		return nil, err
	}

	r.DeletePassword = password
	r.ThreadID = uint(threadId)
	transaction := rs.db.Transaction(func(tx *gorm.DB) error {
		var t models.Thread
		if result := tx.Where("id = ?", tid).First(&t); result.Error != nil {
			return result.Error
		}
		if result := tx.Create(&r); result.Error != nil {
			return result.Error
		}
		if result := tx.Model(&t).Where("id = ?", tid).Update("bumped_on", time.Now()); result != nil {
			return result.Error
		}
		return nil
	})
	return &r, transaction
}

func (rs *replyService) Report(id string) error {
	return rs.db.Transaction(func(tx *gorm.DB) error {
		var r models.Reply
		if result := tx.Where("id = ?", id).First(&r); result.Error != nil {
			return result.Error
		}
		if result := tx.Model(&r).Where("id = ?", id).Update("reported", true); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (rs *replyService) DeleteWithPassword(id, password string) error {
	return rs.db.Transaction(func(tx *gorm.DB) error {
		var r models.Reply
		if result := tx.Where("id = ?", id).First(&r); result.Error != nil {
			return result.Error
		}
		if !utils.CheckPassword(password, r.DeletePassword) {
			return errors.New("incorrect password")
		}
		if result := tx.Model(&r).Where("id = ?", id).Update("reported", true); result.Error != nil {
			return result.Error
		}
		return nil
	})
}
