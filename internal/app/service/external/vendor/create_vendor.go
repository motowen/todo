package externalVendor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
)

type CreateVendorRequest struct {
	VendorName  string `json:"vendor_name" validate:"required" example:"vendor name"`
	VendorAlias string `json:"vendor_alias" example:"vendor alias"`
	Cnty        string `json:"cnty" example:"US"`
}

type CreateVendorResponse struct {
	VendorID string `json:"vendor_id" validate:"required" example:"vendor id"`
}

func (r CreateVendorResponse) Validate() error {
	if r.VendorID == "" {
		return fmt.Errorf("missing vendor_id")
	}
	return nil
}

func CreateVendor(req CreateVendorRequest, token string) (CreateVendorResponse, model.ServiceResp) {
	var response CreateVendorResponse
	api := fmt.Sprintf("%s%s", config.Env.VendorServiceHost, "/api/vendors/v1/vendors")

	httpResp, err := client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetHeader("authKey", token).
		SetBody(req).
		Post(api)

	if err != nil {
		logger.Error.Printf("Failed to CreateVendor: %v", err)
		return response, model.ServiceError.FailedDependencyError(model.ExternalCreateVendorFail)
	}

	if httpResp.StatusCode() != http.StatusOK {
		logger.Error.Printf("CreateVendor httpResp: %v", httpResp)
		return response, model.ServiceError.FailedDependencyError(model.ExternalCreateVendorFail)
	}

	// 手動解析 JSON，以便更好地處理解析錯誤
	if err := json.Unmarshal(httpResp.Body(), &response); err != nil {
		logger.Error.Printf("CreateVendor JSON parse error: %v, body: %s", err, string(httpResp.Body()))
		return response, model.ServiceError.FailedDependencyError(model.ExternalCreateVendorParseFail)
	}

	// 檢查解析結果是否有效 - 這是關鍵的驗證步驟
	// if response.VendorID == "" {
	// 	logger.Error.Printf("CreateVendor response missing vendor_id: %v, body: %s", response, string(httpResp.Body()))
	// 	return response, model.ServiceError.FailedDependencyError(model.ExternalCreateVendorParseFail)
	// }

	// 檢查解析結果是否有效 - 這是關鍵的驗證步驟 更好的做法
	if err := response.Validate(); err != nil {
		logger.Error.Printf("CreateVendor validation failed: %v, body=%s", err, httpResp.String())
		return response, model.ServiceError.FailedDependencyError(model.ExternalCreateVendorParseFail)
	}

	logger.Info.Printf("CreateVendor Success: %v", response)
	return response, model.ServiceError.OK
}
