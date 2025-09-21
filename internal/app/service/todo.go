package service

import (
	"context"
	"net/http"

	externalAccount "go-base/internal/app/service/external/account"
	externalVendor "go-base/internal/app/service/external/vendor"
	"go-base/internal/pkg/database"
	"go-base/internal/pkg/logger"
	"go-base/internal/pkg/model"
	modelDB "go-base/internal/pkg/model/db"
	modelHttp "go-base/internal/pkg/model/http"
	"go-base/internal/pkg/util"
)

// CreateTodo creates a new todo item
func CreateTodo(ctx context.Context, req modelHttp.CreateTodoRequest) (*modelDB.Todo, model.ServiceResp) {
	currentTs := util.GetCurrentMilliseconds()

	token, authServiceResp := externalAccount.GetAuthToken()
	if authServiceResp.Status != http.StatusOK {
		return nil, model.ServiceError.FailedDependencyError(authServiceResp.ErrCode.Code)
	}

	vendorReq := externalVendor.CreateVendorRequest{
		VendorName:  req.Title,
		VendorAlias: req.Description,
		Cnty:        "US",
	}

	vendorResp, createVendorServiceResp := externalVendor.CreateVendor(vendorReq, token)
	if createVendorServiceResp.Status != http.StatusOK {
		return nil, model.ServiceError.FailedDependencyError(createVendorServiceResp.ErrCode.Code)
	}

	todo := &modelDB.Todo{
		ID:          util.GenUUID(),
		Title:       vendorResp.VendorID,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   currentTs,
		UpdatedAt:   currentTs,
	}
	err := database.InsertTodo(*todo)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, model.ServiceError.InternalServiceError(model.DBTimeoutFail)
		}
		return nil, model.ServiceError.InternalServiceError(model.DBCreateTodoFail)
	}

	logger.Info.Printf("Created todo with ID: %s", todo.ID)
	return todo, model.ServiceError.OK
}

func GetAllTodo(ctx context.Context) ([]modelDB.Todo, model.ServiceResp) {
	todos, err := database.GetAllTodo()
	if err != nil {
		return nil, model.ServiceError.InternalServiceError(model.DBFindTodoFail)
	}

	return todos, model.ServiceError.OK
}

func GetTodo(ctx context.Context, id string) (modelDB.Todo, model.ServiceResp) {
	todo, err := database.GetTodo(id)
	if err != nil {
		return modelDB.Todo{}, model.ServiceError.InternalServiceError(model.DBFindTodoFail)
	}

	return todo, model.ServiceError.OK
}

func UpdateTodo(ctx context.Context, id string, req modelHttp.UpdateTodoRequest) model.ServiceResp {
	todo := &modelDB.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		UpdatedAt:   util.GetCurrentMilliseconds(),
	}

	err := database.UpdateTodo(id, *todo)
	if err != nil {
		return model.ServiceError.InternalServiceError(model.DBUpdateTodoFail)
	}

	return model.ServiceError.OK
}

func DeleteTodo(ctx context.Context, id string) model.ServiceResp {
	err := database.DeleteTodo(id)
	if err != nil {
		return model.ServiceError.InternalServiceError(model.DBDeleteTodoFail)
	}

	return model.ServiceError.OK
}
