package externalVendor

import (
	"fmt"
	"net/http"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
)

type GetVendorsResponse struct {
	Vendors []GetVendorResponse `json:"vendors"`
}

func GetVendors(token string) (GetVendorsResponse, model.ServiceResp) {
	var response GetVendorsResponse
	api := fmt.Sprintf("%s%s", config.Env.VendorServiceHost, "/api/vendors/v1/vendors")

	httpResp, err := client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetHeader("authKey", token).
		SetResult(&response).
		Get(api)

	if err != nil {
		logger.Error.Printf("Failed to GetVendors: %v", err)
		return response, model.ServiceError.FailedDependencyError(model.ExternalGetVendorFail)
	}

	if httpResp.StatusCode() != http.StatusOK {
		logger.Error.Printf("GetVendors httpResp: %v", httpResp)
		return response, model.ServiceError.FailedDependencyError(model.ExternalGetVendorFail)
	}

	return response, model.ServiceError.OK
}
