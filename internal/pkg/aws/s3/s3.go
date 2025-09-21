package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"go-base/internal/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	s3SDK "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3API interface {
	PresignPutURL(key string, contentType string) (string, error)
	PresignGetURL(key string) (string, error)
	GetHeadObject(key string) (*s3SDK.HeadObjectOutput, error)
	DeleteObjects(keys []string) (*s3SDK.DeleteObjectsOutput, error)
	CheckObjectExists(key string) (bool, error)
	// New simplified methods
	ListBuckets() ([]types.Bucket, error)
	BucketExists(bucketName string) (bool, error)
	CreateBucket(name string, region string) error
	UploadFile(bucketName string, objectKey string, fileContent []byte) error
	DownloadFile(bucketName string, objectKey string) ([]byte, error)
	CopyObject(sourceBucket string, sourceKey string, destBucket string, destKey string) error
	ListObjects(bucketName string) (interface{}, error)
	DeleteObjectsFromBucket(bucketName string, objectKeys []string) error
	DeleteBucket(bucketName string) error
}

var (
	instance S3API
)

func SetInstance(m S3API) {
	instance = m
}

func GetInstance() S3API {
	return instance
}

type BaseS3API struct {
	client  *s3SDK.Client
	bucket  string
	context context.Context
}

type Config struct {
	AWSS3Bucket         string
	AWSS3Region         string
	IsEnabledAccelerate bool
}

var PresignURLExpiry = 2 * time.Hour

func NewBaseS3API(setupConfig Config) (BaseS3API, error) {
	manager := BaseS3API{}
	background := context.Background()
	cfg, err := config.LoadDefaultConfig(
		background,
		config.WithWebIdentityRoleCredentialOptions(func(o *stscreds.WebIdentityRoleOptions) {
			o.Duration = time.Hour * 12 //max session duration, config by OPS
		}),
		config.WithCredentialsCacheOptions(func(o *aws.CredentialsCacheOptions) {
			o.ExpiryWindow = PresignURLExpiry + time.Minute*1 //this value should greater than PresignURLExpiry, so credential won't expire earlier than url
		}),
	)

	if err != nil {
		logger.Error.Printf("get aws config fail, %+v\n", err)
		return manager, err
	}

	credentials, err := cfg.Credentials.Retrieve(background)
	if err != nil {
		logger.Error.Printf("get aws credentials fail, %+v\n", err)
		return manager, err
	}
	logger.Info.Printf("credentials info: Source = %+v, CanExpire = %+v, Expires = %v\n", credentials.Source, credentials.CanExpire, credentials.Expires)

	client := s3SDK.NewFromConfig(cfg, func(o *s3SDK.Options) {
		o.Region = setupConfig.AWSS3Region
		o.UseAccelerate = setupConfig.IsEnabledAccelerate
	})

	manager = BaseS3API{
		client:  client,
		bucket:  setupConfig.AWSS3Bucket,
		context: background,
	}
	return manager, nil
}

func (manager BaseS3API) PresignPutURL(key string, contentType string) (string, error) {
	psClient := s3SDK.NewPresignClient(manager.client)
	input := &s3SDK.PutObjectInput{
		Bucket:      &manager.bucket,
		Key:         &key,
		ContentType: &contentType,
	}
	resp, err := psClient.PresignPutObject(
		manager.context,
		input,
		s3SDK.WithPresignExpires(PresignURLExpiry),
	)
	if err != nil {
		logger.Error.Printf("PresignPutObject fail, %+v\n", err)
		return "", err
	}
	return resp.URL, err
}

func (manager BaseS3API) PresignGetURL(key string) (string, error) {
	psClient := s3SDK.NewPresignClient(manager.client)
	input := &s3SDK.GetObjectInput{
		Bucket: &manager.bucket,
		Key:    &key,
	}
	resp, err := psClient.PresignGetObject(
		manager.context,
		input,
		s3SDK.WithPresignExpires(PresignURLExpiry),
	)
	if err != nil {
		logger.Error.Printf("PresignGetObject fail, %+v\n", err)
		return "", err
	}
	return resp.URL, err
}

func (manager BaseS3API) GetHeadObject(key string) (*s3SDK.HeadObjectOutput, error) {
	input := &s3SDK.HeadObjectInput{
		Bucket: &manager.bucket,
		Key:    &key,
	}
	return manager.client.HeadObject(manager.context, input)
}

func (manager BaseS3API) DeleteObjects(keys []string) (*s3SDK.DeleteObjectsOutput, error) {
	var objectKeys []types.ObjectIdentifier
	for _, v := range keys {
		objectKeys = append(objectKeys, types.ObjectIdentifier{Key: aws.String(v)})
	}

	input := &s3SDK.DeleteObjectsInput{
		Bucket: &manager.bucket,
		Delete: &types.Delete{
			Objects: objectKeys,
		},
	}

	return manager.client.DeleteObjects(manager.context, input)
}

func (manager BaseS3API) CheckObjectExists(key string) (bool, error) {
	input := &s3SDK.HeadObjectInput{
		Bucket: &manager.bucket,
		Key:    aws.String(key),
	}
	_, err := manager.client.HeadObject(manager.context, input)
	if err != nil {
		var notFoundErr *types.NotFound
		if ok := errors.As(err, &notFoundErr); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// New simplified implementations
func (manager BaseS3API) ListBuckets() ([]types.Bucket, error) {
	/*
		output, err := manager.client.ListBuckets(manager.context, &s3SDK.ListBucketsInput{})
		if err != nil {
			logger.Error.Printf("ListBuckets fail, %+v\n", err)
			return nil, err
		}
		logger.Info.Printf("ListBuckets success, %+v\n", output.Buckets)
		return output.Buckets, nil
	*/
	var err error
	var output *s3SDK.ListBucketsOutput
	var buckets []types.Bucket
	bucketPaginator := s3SDK.NewListBucketsPaginator(manager.client, &s3SDK.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(manager.context)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				logger.Error.Printf("You don't have permission to list buckets for this account.")
				err = apiErr
			} else {
				logger.Error.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}
	return buckets, err
}

func (manager BaseS3API) BucketExists(bucketName string) (bool, error) {
	_, err := manager.client.HeadBucket(manager.context, &s3SDK.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (manager BaseS3API) CreateBucket(name string, region string) error {
	_, err := manager.client.CreateBucket(manager.context, &s3SDK.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	return err
}

func (manager BaseS3API) UploadFile(bucketName string, objectKey string, fileContent []byte) error {
	_, err := manager.client.PutObject(manager.context, &s3SDK.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(fileContent),
		ContentLength: aws.Int64(int64(len(fileContent))),
	})
	return err
}

func (manager BaseS3API) DownloadFile(bucketName string, objectKey string) ([]byte, error) {
	result, err := manager.client.GetObject(manager.context, &s3SDK.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (manager BaseS3API) CopyObject(sourceBucket string, sourceKey string, destBucket string, destKey string) error {
	_, err := manager.client.CopyObject(manager.context, &s3SDK.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		CopySource: aws.String(sourceBucket + "/" + sourceKey),
		Key:        aws.String(destKey),
	})
	return err
}

func (manager BaseS3API) ListObjects(bucketName string) (interface{}, error) {
	output, err := manager.client.ListObjectsV2(manager.context, &s3SDK.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		logger.Error.Printf("ListObjects fail, %+v\n", err)
		return nil, err
	}
	logger.Info.Printf("ListObjects success, %+v\n", output.Contents)
	return output.Contents, nil
}

func (manager BaseS3API) DeleteObjectsFromBucket(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := manager.client.DeleteObjects(manager.context, &s3SDK.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	return err
}

func (manager BaseS3API) DeleteBucket(bucketName string) error {
	_, err := manager.client.DeleteBucket(manager.context, &s3SDK.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}
