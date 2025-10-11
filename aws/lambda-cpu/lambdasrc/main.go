package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CPUInfo struct {
	Id            string `dynamodbav:"id"`
	Timestamp     string `dynamodbav:"timestamp"`
	Architecture  string `dynamodbav:"architecture"`
	MemoryMB      string `dynamodbav:"memory_mb"`
	NumCPU        int    `dynamodbav:"num_cpu"`
	CPUModelName  string `dynamodbav:"cpu_model_name"`
	PhysicalCores int    `dynamodbav:"physical_cores"`
	LogicalCores  int    `dynamodbav:"logical_cores"`
	CPUInfoRaw    string `dynamodbav:"cpu_info_raw"`
}

var (
	dynamoClient *dynamodb.Client
	tableName    = os.Getenv("RESULTS_TABLE_NAME")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

func getCPUInfo(ctx context.Context) (*CPUInfo, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get lambda context")
	}

	info := &CPUInfo{
		Id:           lc.AwsRequestID,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Architecture: runtime.GOARCH,
		MemoryMB:     os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE"),
		NumCPU:       runtime.NumCPU(),
		LogicalCores: runtime.NumCPU(),
	}

	// /proc/cpuinfo から詳細情報を取得
	cpuInfoBytes, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		cpuInfoStr := string(cpuInfoBytes)
		info.CPUInfoRaw = cpuInfoStr

		lines := strings.Split(cpuInfoStr, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					info.CPUModelName = strings.TrimSpace(parts[1])
					break
				}
			}
		}

		// 物理コア数を取得
		coreIDs := make(map[string]bool)
		for _, line := range lines {
			if strings.HasPrefix(line, "core id") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					coreIDs[strings.TrimSpace(parts[1])] = true
				}
			}
		}

		if len(coreIDs) > 0 {
			info.PhysicalCores = len(coreIDs)
		}
	}

	return info, nil
}

func saveToDynamoDB(ctx context.Context, info *CPUInfo) error {
	item, err := attributevalue.MarshalMap(info)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func handler(ctx context.Context) (string, error) {
	info, err := getCPUInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get CPU info: %w", err)
	}

	// DynamoDB に保存
	if err := saveToDynamoDB(ctx, info); err != nil {
		return "", fmt.Errorf("failed to save to DynamoDB: %w", err)
	}

	return fmt.Sprintf("Saved: %s - %s cores", info.CPUModelName, fmt.Sprint(info.LogicalCores)), nil
}

func main() {
	lambda.Start(handler)

}
