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

// 1. ListBuckets Handler
func ListBucketsHandler(c *gin.Context) {
	var request modelHttp.ListBucketsRequest

	ctx := c.Request.Context()
	response, serviceResp := service.ListBuckets(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to list buckets: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 2. BucketExists Handler
func BucketExistsHandler(c *gin.Context) {
	var request modelHttp.BucketExistsRequest

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
	response, serviceResp := service.BucketExists(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to check bucket exists: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 3. CreateBucket Handler
func CreateBucketHandler(c *gin.Context) {
	var request modelHttp.CreateBucketRequest

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
	response, serviceResp := service.CreateBucket(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to create bucket: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 4. UploadFile Handler
func UploadFileHandler(c *gin.Context) {
	var request modelHttp.UploadFileRequest

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
	response, serviceResp := service.UploadFile(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to upload file: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 5. UploadLargeObject Handler
func UploadLargeObjectHandler(c *gin.Context) {
	var request modelHttp.UploadLargeObjectRequest

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
	response, serviceResp := service.UploadLargeObject(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to upload large object: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 6. DownloadFile Handler
func DownloadFileHandler(c *gin.Context) {
	var request modelHttp.DownloadFileRequest

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
	response, serviceResp := service.DownloadFile(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to download file: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 7. DownloadLargeObject Handler
func DownloadLargeObjectHandler(c *gin.Context) {
	var request modelHttp.DownloadLargeObjectRequest

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
	response, serviceResp := service.DownloadLargeObject(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to download large object: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 8. CopyToFolder Handler
func CopyToFolderHandler(c *gin.Context) {
	var request modelHttp.CopyToFolderRequest

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
	response, serviceResp := service.CopyToFolder(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to copy to folder: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 9. CopyToBucket Handler
func CopyToBucketHandler(c *gin.Context) {
	var request modelHttp.CopyToBucketRequest

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
	response, serviceResp := service.CopyToBucket(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to copy to bucket: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 10. ListObjects Handler
func ListObjectsHandler(c *gin.Context) {
	var request modelHttp.ListObjectsRequest

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
	response, serviceResp := service.ListObjects(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to list objects: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 11. DeleteObjectsFromBucket Handler
func DeleteObjectsFromBucketHandler(c *gin.Context) {
	var request modelHttp.DeleteObjectsFromBucketRequest

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
	response, serviceResp := service.DeleteObjectsFromBucket(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to delete objects from bucket: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// 12. DeleteBucket Handler
func DeleteBucketHandler(c *gin.Context) {
	var request modelHttp.DeleteBucketRequest

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
	response, serviceResp := service.DeleteBucket(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to delete bucket: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}
