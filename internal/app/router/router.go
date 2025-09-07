package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/handler"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
)

var Router *gin.Engine

func SetupRouter() (router *gin.Engine) {
	docs.SwaggerInfo.Version = config.Env.Version

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router = gin.Default()
	router.Use(handler.CORSMiddleware(), handler.ErrorMiddleware())

	router.GET("/health", handler.HealthHandler)
	router.GET("/version", handler.VersionHandler)

	// Todo routes
	todoRoutes := router.Group("/todo")
	{
		todoRoutes.GET("", handler.GetAllTodoHandler)
		todoRoutes.GET("/:id", handler.GetTodoHandler)
		todoRoutes.POST("", handler.CreateTodoHandler)
		todoRoutes.PUT("/:id", handler.UpdateTodoHandler)
		todoRoutes.DELETE("/:id", handler.DeleteTodoHandler)
	}

	// S3 routes
	iconRoutes := router.Group("/s3")
	{
		iconRoutes.GET("/presigned-url", handler.GetIconPresignedURLHandler)
		iconRoutes.GET("/head-object", handler.GetIconHeadObjectHandler)
		iconRoutes.GET("/check-object-exists", handler.GetIconCheckObjectExistsHandler)
		iconRoutes.GET("/delete-objects", handler.GetIconDeleteObjectsHandler)

		// New S3 APIs
		iconRoutes.GET("/list-buckets", handler.ListBucketsHandler)
		iconRoutes.GET("/bucket-exists", handler.BucketExistsHandler)
		iconRoutes.POST("/create-bucket", handler.CreateBucketHandler)
		iconRoutes.POST("/upload-file", handler.UploadFileHandler)
		iconRoutes.POST("/upload-large-object", handler.UploadLargeObjectHandler)
		iconRoutes.GET("/download-file", handler.DownloadFileHandler)
		iconRoutes.GET("/download-large-object", handler.DownloadLargeObjectHandler)
		iconRoutes.POST("/copy-to-folder", handler.CopyToFolderHandler)
		iconRoutes.POST("/copy-to-bucket", handler.CopyToBucketHandler)
		iconRoutes.GET("/list-objects", handler.ListObjectsHandler)
		iconRoutes.DELETE("/delete-objects-from-bucket", handler.DeleteObjectsFromBucketHandler)
		iconRoutes.DELETE("/delete-bucket", handler.DeleteBucketHandler)
	}

	// SQS routes
	sqsRoutes := router.Group("/sqs")
	{
		sqsRoutes.POST("/send-message", handler.SendMessageHandler)
		sqsRoutes.POST("/send-messages", handler.SendMessagesHandler)
	}

	return
}

func Setup() error {
	Router = SetupRouter()

	return nil
}
