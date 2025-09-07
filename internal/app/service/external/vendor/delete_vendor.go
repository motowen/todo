package externalVendor

import (
	"fmt"
	"net/http"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
)

type DeleteVendorRequest struct {
	VendorID string `json:"vendor_id" validate:"required" example:"vendor id"`
}

func DeleteVendor(req DeleteVendorRequest, token string) model.ServiceResp {
	api := fmt.Sprintf("%s%s", config.Env.VendorServiceHost, "/api/vendors/v1/vendors?vid=%s", req.VendorID)

	httpResp, err := client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetHeader("authKey", token).
		SetBody(req).
		Delete(api)

	if err != nil {
		logger.Error.Printf("Failed to DeleteVendor: %v", err)
		return model.ServiceError.FailedDependencyError(model.ExternalDeleteVendorFail)
	}

	if httpResp.StatusCode() != http.StatusOK {
		logger.Error.Printf("DeleteVendor httpResp: %v", httpResp)
		return model.ServiceError.FailedDependencyError(model.ExternalDeleteVendorFail)
	}

	return model.ServiceError.OK
}
