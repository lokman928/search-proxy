package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lokman928/search-proxy/internal/module/brave/internal/service"
)

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

type SearchRequestBody struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

func (c *Controller) Search(ctx *gin.Context) {
	var requestBody SearchRequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchResult, err := c.service.Search(ctx, requestBody.Query, requestBody.Count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, searchResult)
}
