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

type ThreadController interface {
	ListThreads(c *gin.Context)
	GetThread(c *gin.Context)
	CreateThread(c *gin.Context)
	Report(c *gin.Context)
	Delete(c *gin.Context)
}

type threadController struct {
	ts services.ThreadService
}

func NewThreadController(ts services.ThreadService) ThreadController {
	return &threadController{ts: ts}
}

func (tc *threadController) CreateThread(c *gin.Context) {
	var t models.Thread
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid request body")))
		return
	}

	thread, err := tc.ts.Create(t)

	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}

	c.JSON(http.StatusCreated, thread)
}

func (tc *threadController) ListThreads(c *gin.Context) {
	page := 0
	pageSize := 5
	pageQuery := c.Query("page")
	pageSizeQuery := c.Query("pageSize")

	if pageQuery != "" {
		p, err := strconv.Atoi(pageQuery)
		if err != nil {
			c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid query parameter")))
			return
		}
		page = p
	}

	if pageSizeQuery != "" {
		p, err := strconv.Atoi(pageSizeQuery)
		if err != nil {
			c.JSON(utils.CreateApiError(http.StatusBadRequest, errors.New("invalid query parameter")))
			return
		}
		pageSize = p
	}

	threads, err := tc.ts.List(page, pageSize)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}
	c.JSON(http.StatusOK, threads)
}

func (tc *threadController) GetThread(c *gin.Context) {
	id := c.Param("id")
	thread, err := tc.ts.GetById(id)

	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (tc *threadController) Report(c *gin.Context) {
	id := c.Param("id")
	err := tc.ts.Report(id)

	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (tc *threadController) Delete(c *gin.Context) {
	id := c.Param("id")
	password := c.Query("password")

	if password == "" {
		c.JSON(utils.ErrorFromDb(errors.New("invalid password")))
		return
	}

	err := tc.ts.DeleteWithPassword(id, password)
	if err != nil {
		c.JSON(utils.ErrorFromDb(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
