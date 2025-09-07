# S3 API 測試 Curl 命令

以下是用於測試 12 個新 S3 API 端點的 curl 命令。請將 `localhost:9999` 替換為您的實際服務器地址。

## 1. ListBuckets - 列出所有存儲桶

```bash
curl -X GET "http://localhost:9999/s3/list-buckets" \
  -H "Content-Type: application/json"
```

## 2. BucketExists - 檢查存儲桶是否存在

```bash
curl -X GET "http://localhost:9999/s3/bucket-exists?bucket_name=htc-enterprise-test-dev" \
  -H "Content-Type: application/json"
```

## 3. CreateBucket - 創建新存儲桶

```bash
curl -X POST "http://localhost:9999/s3/create-bucket" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "htc-enterprise-test-dev-2",
    "region": "us-west-2"
  }'
```

## 4. UploadFile - 上傳文件

```bash
# 首先將文件內容轉換為 base64 編碼
# echo "Hello World!" | base64
curl -X POST "http://localhost:8080/s3/upload-file" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "my-test-bucket",
    "object_key": "test-file.txt",
    "file_data": "SGVsbG8gV29ybGQhCg=="
  }'
```

## 5. UploadLargeObject - 上傳大文件

```bash
curl -X POST "http://localhost:8080/s3/upload-large-object" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "my-test-bucket",
    "object_key": "large-file.bin",
    "file_data": "VGhpcyBpcyBhIGxhcmdlIG9iamVjdCBmb3IgdGVzdGluZw=="
  }'
```

## 6. DownloadFile - 下載文件

```bash
curl -X GET "http://localhost:8080/s3/download-file?bucket_name=my-test-bucket&object_key=test-file.txt" \
  -H "Content-Type: application/json"
```

## 7. DownloadLargeObject - 下載大文件

```bash
curl -X GET "http://localhost:8080/s3/download-large-object?bucket_name=my-test-bucket&object_key=large-file.bin" \
  -H "Content-Type: application/json"
```

## 8. CopyToFolder - 複製對象到文件夾

```bash
curl -X POST "http://localhost:8080/s3/copy-to-folder" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "my-test-bucket",
    "object_key": "test-file.txt",
    "folder_name": "backup"
  }'
```

## 9. CopyToBucket - 複製對象到另一個存儲桶

```bash
curl -X POST "http://localhost:8080/s3/copy-to-bucket" \
  -H "Content-Type: application/json" \
  -d '{
    "source_bucket": "my-test-bucket",
    "destination_bucket": "my-backup-bucket",
    "object_key": "test-file.txt"
  }'
```

## 10. ListObjects - 列出存儲桶中的所有對象

```bash
curl -X GET "http://localhost:9999/s3/list-objects?bucket_name=htc-enterprise-test-dev" \
  -H "Content-Type: application/json"
```

## 11. DeleteObjectsFromBucket - 從存儲桶中刪除多個對象

```bash
curl -X DELETE "http://localhost:8080/s3/delete-objects-from-bucket" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "my-test-bucket",
    "object_keys": ["test-file.txt", "large-file.bin"]
  }'
```

## 12. DeleteBucket - 刪除存儲桶

```bash
curl -X DELETE "http://localhost:8080/s3/delete-bucket" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket_name": "my-old-bucket"
  }'
```

## 其他現有的 S3 API

### GetPresignedURL - 獲取預簽名 URL

```bash
curl -X GET "http://localhost:8080/s3/presigned-url?key=my-file.txt&method=PUT&content_type=text/plain" \
  -H "Content-Type: application/json"
```

### HeadObject - 獲取對象元數據

```bash
curl -X GET "http://localhost:8080/s3/head-object?key=my-file.txt" \
  -H "Content-Type: application/json"
```

### CheckObjectExists - 檢查對象是否存在

```bash
curl -X GET "http://localhost:8080/s3/check-object-exists?key=my-file.txt" \
  -H "Content-Type: application/json"
```

### DeleteObjects - 刪除多個對象

```bash
curl -X GET "http://localhost:8080/s3/delete-objects?keys=file1.txt&keys=file2.txt" \
  -H "Content-Type: application/json"
```

## 注意事項

1. **Base64 編碼**: 上傳文件時，`file_data` 字段需要是 base64 編碼的字符串。
2. **AWS 憑證**: 確保您的服務器配置了正確的 AWS 憑證。
3. **權限**: 確保您有足夠的權限執行這些 S3 操作。
4. **存儲桶名稱**: 存儲桶名稱必須是全局唯一的。
5. **區域**: 創建存儲桶時要指定正確的 AWS 區域。

## 測試流程建議

1. 首先列出現有存儲桶
2. 檢查目標存儲桶是否存在
3. 如果不存在，創建新存儲桶
4. 上傳測試文件
5. 列出存儲桶中的對象
6. 下載文件驗證
7. 複製文件到文件夾或其他存儲桶
8. 清理：刪除對象和存儲桶
