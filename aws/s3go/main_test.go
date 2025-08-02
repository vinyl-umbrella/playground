package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// テスト用の設定を環境変数から取得
func getTestConfig() (string, string, bool) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	objectKey := os.Getenv("S3_OBJECT_KEY")

	if bucketName == "" || objectKey == "" {
		return "", "", false
	}

	return bucketName, objectKey, true
}

// ベンチマーク用のダウンローダーを作成
func createBenchmarkDownloader(partSize int64, concurrency int) (*BucketBasics, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = partSize
		d.Concurrency = concurrency
	})

	return &BucketBasics{
		S3Client:   s3Client,
		Downloader: downloader,
	}, nil
}

// チャンクサイズ別ベンチマーク (チャンクサイズ: [512KB, 1MB, 2MB, 5MB, 10MB], 並行数: 1)
func BenchmarkS3DownloadPartSize(b *testing.B) {
	bucketName, objectKey, ok := getTestConfig()
	if !ok {
		b.Fatal("Failed to get test config. Set S3_BUCKET_NAME and S3_OBJECT_KEY environment variables.")
	}

	testCases := []struct {
		name     string
		partSize int64
	}{
		// {"512KB", 512 * 1024},
		// {"1MB", 1024 * 1024},
		// {"2MB", 2 * 1024 * 1024},
		// {"5MB", 5 * 1024 * 1024},
		// {"10MB", 10 * 1024 * 1024},
		{"20MB", 20 * 1024 * 1024},
		{"50MB", 50 * 1024 * 1024},
		{"100MB", 100 * 1024 * 1024},
		{"200MB", 200 * 1024 * 1024},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			bucketBasics, err := createBenchmarkDownloader(tc.partSize, 16)
			if err != nil {
				b.Fatalf("Failed to create downloader: %v", err)
			}

			tempFile := "/tmp/benchmark_part_" + tc.name
			defer os.Remove(tempFile)

			ctx := context.Background()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				os.Remove(tempFile)

				err := bucketBasics.StreamDL(ctx, bucketName, objectKey, tempFile)
				if err != nil {
					b.Errorf("Download failed: %v", err)
				}
			}
		})
	}
}

// 並行数別ベンチマーク (チャンクサイズ: 1MiB, 並行数: [2^0, 2^1, 2^2, 2^3, 2^4])
func BenchmarkS3DownloadConcurrency(b *testing.B) {
	bucketName, objectKey, ok := getTestConfig()
	if !ok {
		b.Fatal("Failed to get test config. Set S3_BUCKET_NAME and S3_OBJECT_KEY environment variables.")
	}

	testCases := []struct {
		name        string
		concurrency int
	}{
		{"Seq", 1},
		{"Par2", 2},
		{"Par4", 4},
		{"Par8", 8},
		{"Par16", 16},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			bucketBasics, err := createBenchmarkDownloader(1024*1024, tc.concurrency)
			if err != nil {
				b.Fatalf("Failed to create downloader: %v", err)
			}

			tempFile := "/tmp/benchmark_conc_" + tc.name
			defer os.Remove(tempFile)

			ctx := context.Background()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				os.Remove(tempFile)

				err := bucketBasics.StreamDL(ctx, bucketName, objectKey, tempFile)
				if err != nil {
					b.Errorf("Download failed: %v", err)
				}
			}
		})
	}
}

// 速度重視の設定 (チャンクサイズ: 10MiB、並行数: 16)
func BenchmarkS3DownloadSpeedOptimized(b *testing.B) {
	bucketName, objectKey, ok := getTestConfig()
	if !ok {
		b.Fatal("Failed to get test config. Set S3_BUCKET_NAME and S3_OBJECT_KEY environment variables.")
	}

	bucketBasics, err := createBenchmarkDownloader(10*1024*1024, 16)
	if err != nil {
		b.Fatalf("Failed to create downloader: %v", err)
	}

	tempFile := "/tmp/benchmark_speed_optimized"
	defer os.Remove(tempFile)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		os.Remove(tempFile)

		err := bucketBasics.StreamDL(ctx, bucketName, objectKey, tempFile)
		if err != nil {
			b.Errorf("Download failed: %v", err)
		}
	}
}

func BenchmarkSimpleS3GetObject(b *testing.B) {
	bucketName, objectKey, ok := getTestConfig()
	if !ok {
		b.Fatal("Failed to get test config. Set S3_BUCKET_NAME and S3_OBJECT_KEY environment variables.")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		b.Fatalf("Failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	tempFile := "/tmp/benchmark_get_object"
	defer os.Remove(tempFile)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := simpleGetObject(ctx, s3Client, bucketName, objectKey, tempFile)
		if err != nil {
			b.Errorf("GetObject failed: %v", err)
		}
	}
}

// GetObject 後に大きなバッファを確保している悪い例
func BenchmarkBadS3GetObject(b *testing.B) {
	bucketName, objectKey, _ := getTestConfig()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		b.Fatalf("Failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	tempFile := "/tmp/benchmark_get_object"
	defer os.Remove(tempFile)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := badGetObject(ctx, s3Client, bucketName, objectKey, tempFile)
		if err != nil {
			b.Errorf("GetObject failed: %v", err)
		}
	}
}
