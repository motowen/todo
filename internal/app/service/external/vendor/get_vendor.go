package externalVendor

import (
	"fmt"
	"net/http"

	"go-base/internal/pkg/config"
	"go-base/internal/pkg/http/client"
	"go-base/internal/pkg/logger"
	"go-base/internal/pkg/model"
)

type GetVendorRequest struct {
	VendorID string `json:"vendor_id" validate:"required" example:"vendor id"`
}

type GetVendorResponse struct {
	VendorID    string `json:"vendor_id" example:"vendor id"`
	VendorName  string `json:"vendor_name" example:"vendor name"`
	VendorAlias string `json:"vendor_alias" example:"vendor alias"`
	Cnty        string `json:"cnty" example:"US"`
}

func GetVendor(req GetVendorRequest, token string) (GetVendorResponse, model.ServiceResp) {
	var response GetVendorResponse
	api := fmt.Sprintf("%s%s", config.Env.VendorServiceHost, "/api/vendors/v1/vendors/%s", req.VendorID)

	httpResp, err := client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetHeader("authKey", token).
		SetResult(&response).
		Get(api)

	if err != nil {
		logger.Error.Printf("Failed to GetVendor: %v", err)
		return response, model.ServiceError.FailedDependencyError(model.ExternalGetVendorFail)
	}

	if httpResp.StatusCode() != http.StatusOK {
		logger.Error.Printf("GetVendor httpResp: %v", httpResp)
		return response, model.ServiceError.FailedDependencyError(model.ExternalGetVendorFail)
	}

	return response, model.ServiceError.OK
}
