package s3

import (
	"fmt"
)

// Mock S3API 實現，用於測試
type MockS3API struct {
	ShouldFail bool
}

func (m *MockS3API) PresignGetURL(key string) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("mock S3 error")
	}
	return fmt.Sprintf("https://mock-bucket.s3.amazonaws.com/%s?presigned=true", key), nil
}

func (m *MockS3API) PresignPutURL(key string, contentType string) (string, error) {
	if m.ShouldFail {
		return "", fmt.Errorf("mock S3 error")
	}
	return fmt.Sprintf("https://mock-bucket.s3.amazonaws.com/%s?presigned=true&content-type=%s", key, contentType), nil
}
