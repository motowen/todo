package externalVendor

import (
	"fmt"
	"net/http"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
)

type UpdateVendorRequest struct {
	VendorID    string `json:"vendor_id" validate:"required" example:"vendor id"`
	VendorName  string `json:"vendor_name" example:"vendor name"`
	VendorAlias string `json:"vendor_alias" example:"vendor alias"`
	Cnty        string `json:"cnty" example:"US"`
}

type UpdateVendorResponse struct {
	VendorID    string `json:"vendor_id" example:"vendor id"`
	VendorName  string `json:"vendor_name" example:"vendor name"`
	VendorAlias string `json:"vendor_alias" example:"vendor alias"`
	Cnty        string `json:"cnty" example:"US"`
}

func UpdateVendor(req UpdateVendorRequest, token string) (UpdateVendorResponse, model.ServiceResp) {
	var response UpdateVendorResponse
	api := fmt.Sprintf("%s%s", config.Env.VendorServiceHost, "/api/vendors/v1/vendors/%s", req.VendorID)

	httpResp, err := client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetHeader("authKey", token).
		SetBody(req).
		SetResult(&response).
		Put(api)

	if err != nil {
		logger.Error.Printf("Failed to UpdateVendor: %v", err)
		return response, model.ServiceError.FailedDependencyError(model.ExternalUpdateVendorFail)
	}

	if httpResp.StatusCode() != http.StatusOK {
		logger.Error.Printf("UpdateVendor httpResp: %v", httpResp)
		return response, model.ServiceError.FailedDependencyError(model.ExternalUpdateVendorFail)
	}

	return response, model.ServiceError.OK
}
