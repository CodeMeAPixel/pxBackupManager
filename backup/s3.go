package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func UploadToS3(filePath string, bucket string, region string, endpoint string, accessKey string, secretKey string) (string, error) {
	// Create context
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	var client *s3.Client
	if endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	// Extract filename
	filename := filepath.Base(filePath)

	// Upload to S3
	fmt.Printf("Uploading to S3: s3://%s/%s\n", bucket, filename)
	result, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Build S3 URL
	s3URL := fmt.Sprintf("s3://%s/%s", bucket, filename)
	fmt.Printf("Successfully uploaded to S3\n")
	fmt.Printf("ETag: %s\n", aws.ToString(result.ETag))
	fmt.Printf("Size: %.2f MB\n", float64(fileInfo.Size())/1024/1024)

	return s3URL, nil
}

// CheckS3Connection verifies that S3 credentials and connection are valid
func CheckS3Connection(bucket string, region string, endpoint string, accessKey string, secretKey string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	var client *s3.Client
	if endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}

	// Try to list objects (minimal operation to test connection)
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var ae smithy.APIError
		if apierr, ok := err.(smithy.APIError); ok {
			ae = apierr
			return fmt.Errorf("failed to access S3 bucket: %s - %s", ae.ErrorCode(), ae.ErrorMessage())
		}
		return fmt.Errorf("failed to access S3 bucket: %w", err)
	}

	return nil
}
