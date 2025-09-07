package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/router"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/aws/s3"
)

func Test_GetIconPresignedURL_Success_GET(t *testing.T) {
	// 設置 mock S3
	mockS3 := &s3.MockS3API{ShouldFail: false}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png&method=GET", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查響應是否包含 presigned URL
	body := w.Body.String()
	if !strings.Contains(body, "presigned_url") {
		t.Errorf("expected response to contain presigned_url, got: %s", body)
	}
	if !strings.Contains(body, "test/icon.png") {
		t.Errorf("expected response to contain key, got: %s", body)
	}
}

func Test_GetIconPresignedURL_Success_PUT(t *testing.T) {
	// 設置 mock S3
	mockS3 := &s3.MockS3API{ShouldFail: false}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png&method=PUT&content_type=image/png", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查響應是否包含 presigned URL 和 content-type
	body := w.Body.String()
	if !strings.Contains(body, "presigned_url") {
		t.Errorf("expected response to contain presigned_url, got: %s", body)
	}
	if !strings.Contains(body, "test/icon.png") {
		t.Errorf("expected response to contain key, got: %s", body)
	}
	if !strings.Contains(body, "image/png") {
		t.Errorf("expected response to contain content-type, got: %s", body)
	}
}

func Test_GetIconPresignedURL_Invalid_Method(t *testing.T) {
	// 設置 mock S3
	mockS3 := &s3.MockS3API{ShouldFail: false}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png&method=DELETE", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "3001") { // HttpMethodInvalid
		t.Errorf("expected error code 3001, got body: %s", body)
	}
}

func Test_GetIconPresignedURL_Missing_Key(t *testing.T) {
	// 設置 mock S3
	mockS3 := &s3.MockS3API{ShouldFail: false}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?method=GET", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func Test_GetIconPresignedURL_Missing_Method(t *testing.T) {
	// 設置 mock S3
	mockS3 := &s3.MockS3API{ShouldFail: false}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func Test_GetIconPresignedURL_S3_Error(t *testing.T) {
	// 設置 mock S3 為失敗狀態
	mockS3 := &s3.MockS3API{ShouldFail: true}
	s3.SetInstance(mockS3)
	defer s3.SetInstance(nil) // 清理

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png&method=GET", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "1006") { // DBGetIconPresignedURLFail
		t.Errorf("expected error code 1006, got body: %s", body)
	}
}

func Test_GetIconPresignedURL_S3_Not_Initialized(t *testing.T) {
	// 不設置 S3 instance，測試未初始化的情況
	s3.SetInstance(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/icon/presigned-url?key=test/icon.png&method=GET", nil)

	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body=%s", w.Code, w.Body.String())
	}

	// 檢查錯誤碼
	body := w.Body.String()
	if !strings.Contains(body, "1006") { // DBGetIconPresignedURLFail
		t.Errorf("expected error code 1006, got body: %s", body)
	}
}
