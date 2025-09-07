package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/service"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	model "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
	modelHttp "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model/http"
)

// SendMessagesHandler handles sending multiple messages to SQS queue
// @Summary Send multiple messages to SQS queue
// @Description Send multiple messages to specified SQS queue
// @Tags SQS
// @Accept json
// @Produce json
// @Param request body modelHttp.SendMessagesRequest true "Send messages request"
// @Success 200 {object} modelHttp.SendMessagesResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /sqs/send-messages [post]
func SendMessagesHandler(c *gin.Context) {
	var request modelHttp.SendMessagesRequest
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
	response, serviceResp := service.SendMessages(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to send messages: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}

// SendMessageHandler handles sending a single message to SQS queue
// @Summary Send a single message to SQS queue
// @Description Send a single message to specified SQS queue
// @Tags SQS
// @Accept json
// @Produce json
// @Param request body modelHttp.SendMessageRequest true "Send message request"
// @Success 200 {object} modelHttp.SendMessageResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /sqs/send-message [post]
func SendMessageHandler(c *gin.Context) {
	var request modelHttp.SendMessageRequest
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
	response, serviceResp := service.SendMessage(ctx, request)
	if serviceResp.Status != http.StatusOK {
		logger.Error.Printf("Failed to send message: %v", serviceResp.ErrCode)
		result(c, nil, serviceResp)
		return
	}

	result(c, response, serviceResp)
}
