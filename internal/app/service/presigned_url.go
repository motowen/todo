package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/aws/s3"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	model "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
	modelHttp "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model/http"
)

// Http Method
const HttpMethodGet = "GET"
const HttpMethodPut = "PUT"

func GetIconPresignedURL(ctx context.Context, req modelHttp.GetIconPresignedURLRequest) (modelHttp.GetIconPresignedURLResponse, model.ServiceResp) {
	var presignedURL string
	var err error

	if req.Method == HttpMethodGet {
		presignedURL, err = s3.GetInstance().PresignGetURL(req.Key)
		if err != nil {
			return modelHttp.GetIconPresignedURLResponse{}, model.ServiceError.InternalServiceError(model.DBGetIconPresignedURLFail)
		}
	} else if req.Method == HttpMethodPut {
		presignedURL, err = s3.GetInstance().PresignPutURL(req.Key, req.ContentType)
		if err != nil {
			return modelHttp.GetIconPresignedURLResponse{}, model.ServiceError.InternalServiceError(model.DBGetIconPresignedURLFail)
		}
	} else {
		return modelHttp.GetIconPresignedURLResponse{}, model.ServiceError.BadRequestError(model.HttpMethodInvalid)
	}

	response := modelHttp.GetIconPresignedURLResponse{
		PresignedURL: presignedURL,
	}

	return response, model.ServiceError.OK
}

func GetIconHeadObject(ctx context.Context, req modelHttp.GetIconHeadObjectRequest) (modelHttp.GetIconHeadObjectResponse, model.ServiceResp) {
	headObjectOutput, err := s3.GetInstance().GetHeadObject(req.Key)
	logger.Info.Printf("GetIconHeadObject: %+v", headObjectOutput)
	if err != nil {
		// 檢查是否為物件不存在的錯誤
		response := modelHttp.GetIconHeadObjectResponse{
			Exists: false,
		}
		return response, model.ServiceError.OK
	}

	// 構建響應
	response := modelHttp.GetIconHeadObjectResponse{
		Exists: true,
	}

	if headObjectOutput.ContentLength != nil {
		response.ContentLength = *headObjectOutput.ContentLength
	}

	if headObjectOutput.ContentType != nil {
		response.ContentType = *headObjectOutput.ContentType
	}

	if headObjectOutput.LastModified != nil {
		response.LastModified = headObjectOutput.LastModified.Format("2006-01-02T15:04:05Z")
	}

	if headObjectOutput.ETag != nil {
		response.ETag = *headObjectOutput.ETag
	}

	if headObjectOutput.Metadata != nil {
		response.Metadata = headObjectOutput.Metadata
	}

	return response, model.ServiceError.OK
}

func GetIconCheckObjectExists(ctx context.Context, req modelHttp.GetIconCheckObjectExistsRequest) (modelHttp.GetIconCheckObjectExistsResponse, model.ServiceResp) {
	exists, err := s3.GetInstance().CheckObjectExists(req.Key)
	if err != nil {
		return modelHttp.GetIconCheckObjectExistsResponse{}, model.ServiceError.InternalServiceError("Failed to check object existence")
	}

	response := modelHttp.GetIconCheckObjectExistsResponse{
		Exists: exists,
	}

	return response, model.ServiceError.OK
}

func GetIconDeleteObjects(ctx context.Context, req modelHttp.GetIconDeleteObjectsRequest) (modelHttp.GetIconDeleteObjectsResponse, model.ServiceResp) {
	deleteObjectsOutput, err := s3.GetInstance().DeleteObjects(req.Keys)
	if err != nil {
		return modelHttp.GetIconDeleteObjectsResponse{}, model.ServiceError.InternalServiceError("Failed to delete objects")
	}

	response := modelHttp.GetIconDeleteObjectsResponse{
		Success:      true,
		DeletedCount: len(deleteObjectsOutput.Deleted),
	}

	// 處理錯誤項目
	if len(deleteObjectsOutput.Errors) > 0 {
		var errors []modelHttp.DeleteObjectError
		for _, deleteError := range deleteObjectsOutput.Errors {
			errorItem := modelHttp.DeleteObjectError{}

			if deleteError.Key != nil {
				errorItem.Key = *deleteError.Key
			}

			if deleteError.Code != nil {
				errorItem.Code = *deleteError.Code
			}

			if deleteError.Message != nil {
				errorItem.Message = *deleteError.Message
			}

			errors = append(errors, errorItem)
		}
		response.Errors = errors

		// 如果有錯誤，將成功標記為 false
		if len(deleteObjectsOutput.Deleted) == 0 {
			response.Success = false
		}
	}

	return response, model.ServiceError.OK
}

// 1. ListBuckets Service
func ListBuckets(ctx context.Context, req modelHttp.ListBucketsRequest) (modelHttp.ListBucketsResponse, model.ServiceResp) {
	buckets, err := s3.GetInstance().ListBuckets()
	if err != nil {
		return modelHttp.ListBucketsResponse{}, model.ServiceError.InternalServiceError("Failed to list buckets")
	}

	response := ConvertBuckets(buckets)
	logger.Info.Printf("ListBuckets response: %+v", response)
	return response, model.ServiceError.OK
}

func ConvertBuckets(buckets []types.Bucket) modelHttp.ListBucketsResponse {
	var resp modelHttp.ListBucketsResponse
	for _, b := range buckets {
		resp.Buckets = append(resp.Buckets, modelHttp.BucketInfo{
			Name:         safeString(b.Name),
			CreationDate: safeTime(b.CreationDate),
		})
	}
	return resp
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339) // or other format like "2006-01-02 15:04:05"
}

// 2. BucketExists Service
func BucketExists(ctx context.Context, req modelHttp.BucketExistsRequest) (modelHttp.BucketExistsResponse, model.ServiceResp) {
	exists, err := s3.GetInstance().BucketExists(req.BucketName)
	if err != nil {
		return modelHttp.BucketExistsResponse{}, model.ServiceError.InternalServiceError("Failed to check bucket existence")
	}

	response := modelHttp.BucketExistsResponse{
		Exists: exists,
	}

	return response, model.ServiceError.OK
}

// 3. CreateBucket Service
func CreateBucket(ctx context.Context, req modelHttp.CreateBucketRequest) (modelHttp.CreateBucketResponse, model.ServiceResp) {
	err := s3.GetInstance().CreateBucket(req.BucketName, req.Region)
	if err != nil {
		return modelHttp.CreateBucketResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create bucket: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to create bucket")
	}

	response := modelHttp.CreateBucketResponse{
		Success: true,
		Message: "Bucket created successfully",
	}

	return response, model.ServiceError.OK
}

// 4. UploadFile Service
func UploadFile(ctx context.Context, req modelHttp.UploadFileRequest) (modelHttp.UploadFileResponse, model.ServiceResp) {
	// Decode base64 file data
	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		return modelHttp.UploadFileResponse{
			Success: false,
			Message: "Invalid file data encoding",
		}, model.ServiceError.BadRequestError("Invalid file data encoding")
	}

	err = s3.GetInstance().UploadFile(req.BucketName, req.ObjectKey, fileData)
	if err != nil {
		return modelHttp.UploadFileResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload file: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to upload file")
	}

	response := modelHttp.UploadFileResponse{
		Success: true,
		Message: "File uploaded successfully",
	}

	return response, model.ServiceError.OK
}

// 5. UploadLargeObject Service
func UploadLargeObject(ctx context.Context, req modelHttp.UploadLargeObjectRequest) (modelHttp.UploadLargeObjectResponse, model.ServiceResp) {
	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		return modelHttp.UploadLargeObjectResponse{
			Success: false,
			Message: "Invalid file data encoding",
		}, model.ServiceError.BadRequestError("Invalid file data encoding")
	}

	err = s3.GetInstance().UploadFile(req.BucketName, req.ObjectKey, fileData)
	if err != nil {
		return modelHttp.UploadLargeObjectResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload large object: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to upload large object")
	}

	response := modelHttp.UploadLargeObjectResponse{
		Success: true,
		Message: "Large object uploaded successfully",
	}

	return response, model.ServiceError.OK
}

// 6. DownloadFile Service
func DownloadFile(ctx context.Context, req modelHttp.DownloadFileRequest) (modelHttp.DownloadFileResponse, model.ServiceResp) {
	fileData, err := s3.GetInstance().DownloadFile(req.BucketName, req.ObjectKey)
	if err != nil {
		return modelHttp.DownloadFileResponse{}, model.ServiceError.InternalServiceError("Failed to download file")
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)

	response := modelHttp.DownloadFileResponse{
		FileData:    encodedData,
		ContentType: "application/octet-stream",
		Size:        int64(len(fileData)),
	}

	return response, model.ServiceError.OK
}

// 7. DownloadLargeObject Service
func DownloadLargeObject(ctx context.Context, req modelHttp.DownloadLargeObjectRequest) (modelHttp.DownloadLargeObjectResponse, model.ServiceResp) {
	fileData, err := s3.GetInstance().DownloadFile(req.BucketName, req.ObjectKey)
	if err != nil {
		return modelHttp.DownloadLargeObjectResponse{}, model.ServiceError.InternalServiceError("Failed to download large object")
	}

	encodedData := base64.StdEncoding.EncodeToString(fileData)

	response := modelHttp.DownloadLargeObjectResponse{
		FileData:    encodedData,
		ContentType: "application/octet-stream",
		Size:        int64(len(fileData)),
	}

	return response, model.ServiceError.OK
}

// 8. CopyToFolder Service
func CopyToFolder(ctx context.Context, req modelHttp.CopyToFolderRequest) (modelHttp.CopyToFolderResponse, model.ServiceResp) {
	newObjectKey := fmt.Sprintf("%s/%s", req.FolderName, req.ObjectKey)

	err := s3.GetInstance().CopyObject(req.BucketName, req.ObjectKey, req.BucketName, newObjectKey)
	if err != nil {
		return modelHttp.CopyToFolderResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to copy to folder: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to copy object to folder")
	}

	response := modelHttp.CopyToFolderResponse{
		Success:      true,
		NewObjectKey: newObjectKey,
		Message:      "Object copied to folder successfully",
	}

	return response, model.ServiceError.OK
}

// 9. CopyToBucket Service
func CopyToBucket(ctx context.Context, req modelHttp.CopyToBucketRequest) (modelHttp.CopyToBucketResponse, model.ServiceResp) {
	err := s3.GetInstance().CopyObject(req.SourceBucket, req.ObjectKey, req.DestinationBucket, req.ObjectKey)
	if err != nil {
		return modelHttp.CopyToBucketResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to copy to bucket: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to copy object to bucket")
	}

	response := modelHttp.CopyToBucketResponse{
		Success: true,
		Message: "Object copied to bucket successfully",
	}

	return response, model.ServiceError.OK
}

// 10. ListObjects Service
func ListObjects(ctx context.Context, req modelHttp.ListObjectsRequest) (modelHttp.ListObjectsResponse, model.ServiceResp) {
	_, err := s3.GetInstance().ListObjects(req.BucketName)
	if err != nil {
		return modelHttp.ListObjectsResponse{}, model.ServiceError.InternalServiceError("Failed to list objects")
	}

	var objects []modelHttp.ObjectInfo
	response := modelHttp.ListObjectsResponse{
		Objects: objects,
	}

	return response, model.ServiceError.OK
}

// 11. DeleteObjectsFromBucket Service
func DeleteObjectsFromBucket(ctx context.Context, req modelHttp.DeleteObjectsFromBucketRequest) (modelHttp.DeleteObjectsFromBucketResponse, model.ServiceResp) {
	err := s3.GetInstance().DeleteObjectsFromBucket(req.BucketName, req.ObjectKeys)
	if err != nil {
		return modelHttp.DeleteObjectsFromBucketResponse{
			Success: false,
		}, model.ServiceError.InternalServiceError("Failed to delete objects from bucket")
	}

	response := modelHttp.DeleteObjectsFromBucketResponse{
		Success:      true,
		DeletedCount: len(req.ObjectKeys),
	}

	return response, model.ServiceError.OK
}

// 12. DeleteBucket Service
func DeleteBucket(ctx context.Context, req modelHttp.DeleteBucketRequest) (modelHttp.DeleteBucketResponse, model.ServiceResp) {
	err := s3.GetInstance().DeleteBucket(req.BucketName)
	if err != nil {
		return modelHttp.DeleteBucketResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete bucket: %v", err),
		}, model.ServiceError.InternalServiceError("Failed to delete bucket")
	}

	response := modelHttp.DeleteBucketResponse{
		Success: true,
		Message: "Bucket deleted successfully",
	}

	return response, model.ServiceError.OK
}
