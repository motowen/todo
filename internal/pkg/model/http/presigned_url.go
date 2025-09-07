package model

type GetIconPresignedURLRequest struct {
	Key         string `json:"key" form:"key" binding:"required"`
	Method      string `json:"method" form:"method" binding:"required"`
	ContentType string `json:"content_type" form:"content_type"`
}

type GetIconPresignedURLResponse struct {
	PresignedURL string `json:"presigned_url"`
}

// Head Object Request and Response
type GetIconHeadObjectRequest struct {
	Key string `json:"key" form:"key" binding:"required"`
}

type GetIconHeadObjectResponse struct {
	Exists        bool              `json:"exists"`
	ContentType   string            `json:"content_type,omitempty"`
	ContentLength int64             `json:"content_length,omitempty"`
	LastModified  string            `json:"last_modified,omitempty"`
	ETag          string            `json:"etag,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// Check Object Exists Request and Response
type GetIconCheckObjectExistsRequest struct {
	Key string `json:"key" form:"key" binding:"required"`
}

type GetIconCheckObjectExistsResponse struct {
	Exists bool `json:"exists"`
}

// Delete Objects Request and Response
type GetIconDeleteObjectsRequest struct {
	Keys []string `json:"keys" form:"keys" binding:"required"`
}

type GetIconDeleteObjectsResponse struct {
	Success      bool                `json:"success"`
	DeletedCount int                 `json:"deleted_count"`
	Errors       []DeleteObjectError `json:"errors,omitempty"`
}

type DeleteObjectError struct {
	Key     string `json:"key"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 1. ListBuckets Request and Response
type ListBucketsRequest struct {
	// No parameters needed
}

type ListBucketsResponse struct {
	Buckets []BucketInfo `json:"buckets"`
}

type BucketInfo struct {
	Name         string `json:"name"`
	CreationDate string `json:"creation_date"`
}

// 2. BucketExists Request and Response
type BucketExistsRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
}

type BucketExistsResponse struct {
	Exists bool `json:"exists"`
}

// 3. CreateBucket Request and Response
type CreateBucketRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	Region     string `json:"region" form:"region" binding:"required"`
}

type CreateBucketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 4. UploadFile Request and Response
type UploadFileRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKey  string `json:"object_key" form:"object_key" binding:"required"`
	FileData   string `json:"file_data" form:"file_data" binding:"required"` // base64 encoded
}

type UploadFileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 5. UploadLargeObject Request and Response
type UploadLargeObjectRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKey  string `json:"object_key" form:"object_key" binding:"required"`
	FileData   string `json:"file_data" form:"file_data" binding:"required"` // base64 encoded
}

type UploadLargeObjectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 6. DownloadFile Request and Response
type DownloadFileRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKey  string `json:"object_key" form:"object_key" binding:"required"`
}

type DownloadFileResponse struct {
	FileData    string `json:"file_data"` // base64 encoded
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// 7. DownloadLargeObject Request and Response
type DownloadLargeObjectRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKey  string `json:"object_key" form:"object_key" binding:"required"`
}

type DownloadLargeObjectResponse struct {
	FileData    string `json:"file_data"` // base64 encoded
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// 8. CopyToFolder Request and Response
type CopyToFolderRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKey  string `json:"object_key" form:"object_key" binding:"required"`
	FolderName string `json:"folder_name" form:"folder_name" binding:"required"`
}

type CopyToFolderResponse struct {
	Success      bool   `json:"success"`
	NewObjectKey string `json:"new_object_key"`
	Message      string `json:"message"`
}

// 9. CopyToBucket Request and Response
type CopyToBucketRequest struct {
	SourceBucket      string `json:"source_bucket" form:"source_bucket" binding:"required"`
	DestinationBucket string `json:"destination_bucket" form:"destination_bucket" binding:"required"`
	ObjectKey         string `json:"object_key" form:"object_key" binding:"required"`
}

type CopyToBucketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// 10. ListObjects Request and Response
type ListObjectsRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
}

type ListObjectsResponse struct {
	Objects []ObjectInfo `json:"objects"`
}

type ObjectInfo struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
	ETag         string `json:"etag"`
}

// 11. DeleteObjectsFromBucket Request and Response
type DeleteObjectsFromBucketRequest struct {
	BucketName string   `json:"bucket_name" form:"bucket_name" binding:"required"`
	ObjectKeys []string `json:"object_keys" form:"object_keys" binding:"required"`
}

type DeleteObjectsFromBucketResponse struct {
	Success      bool                `json:"success"`
	DeletedCount int                 `json:"deleted_count"`
	Errors       []DeleteObjectError `json:"errors,omitempty"`
}

// 12. DeleteBucket Request and Response
type DeleteBucketRequest struct {
	BucketName string `json:"bucket_name" form:"bucket_name" binding:"required"`
}

type DeleteBucketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
