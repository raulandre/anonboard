package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raulandre/anonboard/models"
	"github.com/raulandre/anonboard/services"
	"github.com/raulandre/anonboard/utils"
)

type ReplyController interface {
	ListReplies(c *gin.Context)
	GetReply(c *gin.Context)
	CreateReply(c *gin.Context)
	ReportReply(c *gin.Context)
	DeleteReply(c *gin.Context)
}

type replyController struct {
	rs services.ReplyService
}

func NewReplyController(rs services.ReplyService) ReplyController {
	return &replyController{rs}
}

func (rc *replyController) ListReplies(c *gin.Context) {
	tid := c.Param("tid")
	page := 0
	pageQuery := c.Query("page")
	if pageQuery != "" {
		p, err := strconv.Atoi(pageQuery)
		if err != nil {
			c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid page query parameter")))
			return
		}
		page = p
	}

	threadId, error := strconv.Atoi(tid)
	if error != nil {
		c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid thread id")))
	}

	replies, err := rc.rs.ListByThreadId(threadId, page, 5)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusOK, replies)
}

func (rc *replyController) GetReply(c *gin.Context) {
	id := c.Param("id")
	reply, err := rc.rs.GetById(id)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusOK, reply)
}

func (rc *replyController) CreateReply(c *gin.Context) {
	tid := c.Param("tid")
	var r models.Reply
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid request body")))
	}
	reply, err := rc.rs.Create(tid, r)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusCreated, reply)
}

func (rc *replyController) ReportReply(c *gin.Context) {
	id := c.Param("id")
	err := rc.rs.Report(id)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (rc *replyController) DeleteReply(c *gin.Context) {
	id := c.Param("id")
	password := c.Query("password")
	if password == "" {
		c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("must provide password query")))
		return
	}
	err := rc.rs.DeleteWithPassword(id, password)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
