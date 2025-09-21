# AWS S3 Integration

## 概述

這個模組提供了 AWS S3 的 presigned URL 功能，支援 GET 和 PUT 操作。

## 當前狀態

- ✅ **接口定義**: 完整的 S3API 接口
- ✅ **Mock 實現**: 用於測試的 MockS3API
- ✅ **便利函數**: PresignGetURL() 和 PresignPutURL()
- 🚧 **AWS SDK 實現**: 已實現但暫時註釋掉（避免依賴問題）

## 生產環境啟用步驟

### 1. 安裝 AWS SDK 依賴

```bash
go get github.com/aws/aws-sdk-go-v2/aws@latest
go get github.com/aws/aws-sdk-go-v2/config@latest
go get github.com/aws/aws-sdk-go-v2/credentials@latest
go get github.com/aws/aws-sdk-go-v2/service/s3@latest
```

### 2. 解除註釋 AWS SDK 代碼

在 `internal/pkg/aws/s3/s3.go` 中：

1. 解除註釋 import 區塊中的 AWS SDK imports
2. 解除註釋 `BaseS3API` 結構體和相關函數
3. 如果需要，可以解除註釋其他 S3 操作函數

### 3. 設置 AWS 認證

確保您的環境已配置 AWS 認證，可以通過：
- AWS IAM 角色（推薦用於 EC2/EKS）
- AWS CLI 配置 (`~/.aws/credentials`)
- 環境變數 (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)

### 4. 配置 S3 實例

在您的應用程式啟動時：

```go
import "go-base/internal/pkg/aws/s3"

// 設置 S3 配置
s3Config := s3.Config{
    AWSS3Bucket:         "your-bucket-name",
    AWSS3Region:         "us-west-2",
    IsEnabledAccelerate: false,
}

// 創建 S3 實例
s3Instance, err := s3.NewBaseS3API(s3Config)
if err != nil {
    log.Fatal("Failed to initialize S3:", err)
}

// 設置全局實例
s3.SetInstance(s3Instance)
```

### 5. 添加配置變數

在 `internal/pkg/config/config.go` 中添加：

```go
type EnvVars struct {
    // ... 現有配置 ...
    
    // S3 配置
    AWSS3Bucket         string `env:"AWS_S3_BUCKET,required"`
    AWSS3Region         string `env:"AWS_S3_REGION" envDefault:"us-west-2"`
    IsEnabledAccelerate bool   `env:"AWS_S3_ACCELERATE" envDefault:"false"`
}
```

### 6. 更新環境變數

在 `env.sh` 中添加：

```bash
export AWS_S3_BUCKET='your-bucket-name'
export AWS_S3_REGION='us-west-2'
export AWS_S3_ACCELERATE='false'
```

## API 使用方式

### 獲取 GET presigned URL

```bash
GET /icon/presigned-url?key=icons/my-icon.png&method=GET
```

### 獲取 PUT presigned URL

```bash
GET /icon/presigned-url?key=icons/my-icon.png&method=PUT&content_type=image/png
```

## 測試

所有測試都使用 MockS3API，無需實際的 AWS 認證：

```bash
go test ./test -v -run Test_GetIconPresignedURL
```

## 安全注意事項

1. **權限控制**: 確保 IAM 角色只有必要的 S3 權限
2. **桶策略**: 配置適當的 S3 桶策略
3. **網路安全**: 在生產環境中使用 HTTPS
4. **日誌記錄**: 監控 presigned URL 的使用情況

## 故障排除

### 常見問題

1. **"S3 instance not initialized"**: 確保在應用程式啟動時調用了 `s3.SetInstance()`
2. **AWS 認證錯誤**: 檢查 AWS 認證配置
3. **權限錯誤**: 確保 IAM 角色有 `s3:GetObject` 和 `s3:PutObject` 權限

### 日誌檢查

查看應用程式日誌中的 S3 相關錯誤：

```bash
grep -i "s3\|presigned" /path/to/app.log
```

