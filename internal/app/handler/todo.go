package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	"go-base/internal/app/service"
	"go-base/internal/pkg/logger"
	model "go-base/internal/pkg/model"
	modelHttp "go-base/internal/pkg/model/http"
)

func CreateTodoHandler(c *gin.Context) {
	var request modelHttp.CreateTodoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error.Printf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		logger.Error.Printf("Failed to validate request: %v", err)
		result(c, nil, model.ServiceError.BadRequestError("Validation failed: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	todo, serviceResp := service.CreateTodo(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to create todo: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, todo, serviceResp)
}

func GetAllTodoHandler(c *gin.Context) {
	ctx := c.Request.Context()
	todo, serviceResp := service.GetAllTodo(ctx)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to get all todo: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, todo, serviceResp)
}

func GetTodoHandler(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	todo, serviceResp := service.GetTodo(ctx, id)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to get todo: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, todo, serviceResp)
}

func UpdateTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var request modelHttp.UpdateTodoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error.Printf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		logger.Error.Printf("Failed to validate request: %v", err)
		result(c, nil, model.ServiceError.BadRequestError("Validation failed: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	serviceResp := service.UpdateTodo(ctx, id, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to update todo: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, nil, serviceResp)
}

func DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	serviceResp := service.DeleteTodo(ctx, id)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to delete todo: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, nil, serviceResp)
}
