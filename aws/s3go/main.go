package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client   *s3.Client
	Downloader *manager.Downloader
}

func (b *BucketBasics) StreamDL(ctx context.Context, bucketName, objectKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// stream download
	_, err = b.Downloader.Download(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to download object: %w", err)
	}

	if fileInfo, err := os.Stat(filename); err == nil {
		log.Printf("Downloaded %d bytes to %s", fileInfo.Size(), filename)
	}

	return nil
}

func simpleGetObject(ctx context.Context, s3Client *s3.Client, bucketName, objectKey, filename string) error {
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to download object: %w", err)
	}
	defer result.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("failed to copy object data: %w", err)
	}

	return nil
}

// 単純な S3 GetObject を実行する関数
func badGetObject(ctx context.Context, s3Client *s3.Client, bucketName, objectKey, filename string) error {
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to download object: %w", err)
	}
	defer result.Body.Close()

	fileSize := result.ContentLength
	buffer := make([]byte, *fileSize) // とっても悪い例

	totalRead := 0
	for {
		n, err := result.Body.Read(buffer[totalRead:])
		totalRead += n
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(buffer[:totalRead])
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func main() {
	// load bucket, key from argv
	if len(os.Args) < 4 {
		log.Fatalf("Usage: %s <bucket> <key> <download_file>", os.Args[0])
	}
	bucketName := os.Args[1]
	objectKey := os.Args[2]
	downloadFile := os.Args[3]

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = false
	})
	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = 10 * 1024 * 1024 // 10 MiB
		d.Concurrency = 4             // 4並列
	})

	bucketBasics := &BucketBasics{
		S3Client:   s3Client,
		Downloader: downloader,
	}

	ctx := context.Background()

	log.Printf("Starting download of %s/%s to %s", bucketName, objectKey, downloadFile)

	// ファイルにダウンロード
	err = bucketBasics.StreamDL(ctx, bucketName, objectKey, downloadFile)
	if err != nil {
		log.Printf("Failed to download file: %v\n", err)
		return
	}

	log.Printf("Successfully downloaded file: %s\n", downloadFile)

	// ダウンロードしたファイルの情報を表示
	if fileInfo, err := os.Stat(downloadFile); err == nil {
		log.Printf("Downloaded file size: %.2f MB", float64(fileInfo.Size())/(1024*1024))
	}
}
