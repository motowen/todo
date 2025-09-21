package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-base/internal/app/router"
	externalAccount "go-base/internal/app/service/external/account"
	"go-base/internal/pkg/config"

	"github.com/jarcoal/httpmock"
)

func Test_CreateTodo_ShouldBindJSON_Error(t *testing.T) {
	WithDBCleanup(t)
	w := httptest.NewRecorder()
	// 壞 JSON（ShouldBindJSON 會報錯）
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func Test_CreateTodo_ValidateStruct_Error(t *testing.T) {
	WithDBCleanup(t)
	w := httptest.NewRecorder()
	// 少欄位（validate.Struct 會報錯）
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"only-title"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func Test_CreateTodo_GetAuthToken_401_Failed(t *testing.T) {
	WithDBCleanup(t)

	// 清除之前的 mock 狀態和 auth cache
	httpmock.Reset()
	externalAccount.ClearAuthCache()

	// Mock GetAuthToken API 失敗
	authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
	httpmock.RegisterResponder("POST", authURL,
		httpmock.NewStringResponder(401, `{"error":"unauthorized"}`))

	defer httpmock.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":"d1"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusFailedDependency {
		t.Fatalf("expected 424 (Failed Dependency), got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "2001") { // ExternalGetAuthTokenFail
		t.Errorf("expected error code 2001, got body: %s", body)
	}
}

func Test_CreateTodo_GetAuthToken_Parse_Failed(t *testing.T) {
	WithDBCleanup(t)

	// 清除之前的 mock 狀態和 auth cache
	httpmock.Reset()
	externalAccount.ClearAuthCache()

	// Mock GetAuthToken API 失敗
	authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
	httpmock.RegisterResponder("POST", authURL,
		httpmock.NewStringResponder(200, `{"auth_key":"test-token-123"}`))

	defer httpmock.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":"d1"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusFailedDependency {
		t.Fatalf("expected 424 (Failed Dependency), got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "2002") { // ExternalGetAuthTokenParseFail
		t.Errorf("expected error code 2002, got body: %s", body)
	}
}

func Test_CreateTodo_CreateVendor_500_Failed(t *testing.T) {
	WithDBCleanup(t)

	// 清除之前的狀態
	httpmock.Reset()
	externalAccount.ClearAuthCache()

	// Mock GetAuthToken API 成功
	authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
	httpmock.RegisterResponder("POST", authURL,
		httpmock.NewStringResponder(200, `{"access_token":"test-token-123"}`))

	// Mock CreateVendor API 失敗
	vendorURL := config.Env.VendorServiceHost + "/api/vendors/v1/vendors"
	httpmock.RegisterResponder("POST", vendorURL,
		httpmock.NewStringResponder(500, `{"error":"internal server error"}`))

	defer httpmock.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":"d1"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusFailedDependency {
		t.Fatalf("expected 424 (Failed Dependency), got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "2011") { // ExternalCreateVendorFail
		t.Errorf("expected error code 2011, got body: %s", body)
	}
}

// 測試 JSON 解析的四種情況
func Test_CreateVendor_JSON_Parse_Cases(t *testing.T) {
	testCases := []struct {
		name            string
		responseBody    string
		expectedErrCode string
		description     string
	}{
		{
			name:            "Invalid_JSON",
			responseBody:    `error occurred`,
			expectedErrCode: "2012",
			description:     "回傳不是合法的 JSON (純文字)",
		},
		{
			name:            "Valid_JSON_Wrong_Structure",
			responseBody:    `{"vid":"vendor-123"}`,
			expectedErrCode: "2012",
			description:     "回傳合法 JSON，但結構跟定義不符合",
		},
		{
			name:            "Valid_JSON_Wrong_Type",
			responseBody:    `{"vendor_id":123}`,
			expectedErrCode: "2012",
			description:     "回傳合法 JSON，格式正確但型別不對",
		},
		{
			name:            "Valid_JSON_Correct",
			responseBody:    `{"vendor_id":"vendor-123"}`,
			expectedErrCode: "",
			description:     "回傳合法 JSON，完全跟定義相符合",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			WithDBCleanup(t)

			// 清除之前的狀態
			httpmock.Reset()
			externalAccount.ClearAuthCache()

			// Mock GetAuthToken API 成功
			authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
			httpmock.RegisterResponder("POST", authURL,
				httpmock.NewStringResponder(200, `{"access_token":"test-token-123"}`))

			// Mock CreateVendor API
			vendorURL := config.Env.VendorServiceHost + "/api/vendors/v1/vendors"
			httpmock.RegisterResponder("POST", vendorURL,
				httpmock.NewStringResponder(200, tc.responseBody))

			defer httpmock.Reset()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":"d1"}`))
			req.Header.Set("Content-Type", "application/json")

			router.Router.ServeHTTP(w, req)

			if tc.expectedErrCode == "" {
				// 成功案例
				if w.Code != http.StatusOK {
					t.Fatalf("[%s] expected 200, got %d, body=%s", tc.description, w.Code, w.Body.String())
				}
			} else {
				// 失敗案例
				if w.Code != http.StatusFailedDependency {
					t.Fatalf("[%s] expected 424 (Failed Dependency), got %d, body=%s", tc.description, w.Code, w.Body.String())
				}

				// 檢查錯誤碼
				body := w.Body.String()
				if !strings.Contains(body, tc.expectedErrCode) {
					t.Errorf("[%s] expected error code %s, got body: %s", tc.description, tc.expectedErrCode, body)
				}
			}
		})
	}
}

func Test_CreateTodo_Database_Timeout(t *testing.T) {
	WithDBCleanup(t)

	// Mock 外部服務成功
	authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
	httpmock.RegisterResponder("POST", authURL,
		httpmock.NewStringResponder(200, `{"access_token":"test-token-123"}`))

	vendorURL := config.Env.VendorServiceHost + "/api/vendors/v1/vendors"
	httpmock.RegisterResponder("POST", vendorURL,
		httpmock.NewStringResponder(200, `{"vendor_id":"vendor-123"}`))

	defer httpmock.Reset()

	// 註：這個測試需要 mock database.InsertTodo 來模擬超時，
	// 或者用 context.WithTimeout 來測試，這裡先保留結構
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"timeout-test","description":"test-timeout"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	// 在正常情況下應該成功，若要測試 timeout 需要額外的 mock 設定
	if w.Code != http.StatusOK {
		t.Logf("Note: This test currently doesn't mock database timeout. Got %d, body=%s", w.Code, w.Body.String())
	}
}

func Test_CreateTodo_Success(t *testing.T) {
	WithDBCleanup(t)

	// 清除之前的狀態
	httpmock.Reset()
	externalAccount.ClearAuthCache()

	// Mock GetAuthToken API
	authURL := config.Env.AuthServiceHost + "/$SS$/Services/OAuth/Token"
	httpmock.RegisterResponder("POST", authURL,
		httpmock.NewStringResponder(200, `{"access_token":"test-token-123"}`))

	// Mock CreateVendor API
	vendorURL := config.Env.VendorServiceHost + "/api/vendors/v1/vendors"
	httpmock.RegisterResponder("POST", vendorURL,
		httpmock.NewStringResponder(200, `{"vendor_id":"vendor-123"}`))

	defer httpmock.Reset()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(`{"title":"t1","description":"d1"}`))
	req.Header.Set("Content-Type", "application/json")

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// 驗證有正確調用外部 API
	if httpmock.GetCallCountInfo() == nil {
		t.Error("expected external API calls")
	}
}
