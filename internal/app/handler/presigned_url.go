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

func GetIconPresignedURLHandler(c *gin.Context) {
	var request modelHttp.GetIconPresignedURLRequest

	// 使用查詢參數而不是 JSON body
	if err := c.ShouldBindQuery(&request); err != nil {
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
	response, serviceResp := service.GetIconPresignedURL(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to get icon presigned url: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

func GetIconHeadObjectHandler(c *gin.Context) {
	var request modelHttp.GetIconHeadObjectRequest

	// 使用查詢參數
	if err := c.ShouldBindQuery(&request); err != nil {
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
	response, serviceResp := service.GetIconHeadObject(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to get icon head object: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

func GetIconCheckObjectExistsHandler(c *gin.Context) {
	var request modelHttp.GetIconCheckObjectExistsRequest

	// 使用查詢參數
	if err := c.ShouldBindQuery(&request); err != nil {
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
	response, serviceResp := service.GetIconCheckObjectExists(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to check object exists: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

func GetIconDeleteObjectsHandler(c *gin.Context) {
	var request modelHttp.GetIconDeleteObjectsRequest

	// 使用查詢參數
	if err := c.ShouldBindQuery(&request); err != nil {
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
	response, serviceResp := service.GetIconDeleteObjects(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to delete objects: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}
