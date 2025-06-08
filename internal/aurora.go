package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// StartAuroraCluster Auroraクラスタを起動する
func StartAuroraCluster(awsCtx AwsContext, clusterId string) error {
	cfg, err := LoadAwsConfig(awsCtx)
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := rds.NewFromConfig(cfg)
	_, err = client.StartDBCluster(context.Background(), &rds.StartDBClusterInput{
		DBClusterIdentifier: aws.String(clusterId),
	})
	if err != nil {
		return fmt.Errorf("❌ Aurora DBクラスターの起動に失敗: %w", err)
	}
	return nil
}

// StopAuroraCluster Auroraクラスタを停止する
func StopAuroraCluster(awsCtx AwsContext, clusterId string) error {
	cfg, err := LoadAwsConfig(awsCtx)
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := rds.NewFromConfig(cfg)
	_, err = client.StopDBCluster(context.Background(), &rds.StopDBClusterInput{
		DBClusterIdentifier: aws.String(clusterId),
	})
	if err != nil {
		return fmt.Errorf("❌ Aurora DBクラスターの停止に失敗: %w", err)
	}
	return nil
}

// GetAuroraFromStack はCloudFormationスタック名からAuroraクラスター識別子を取得します。
func GetAuroraFromStack(awsCtx AwsContext, stackName string) (string, error) {
	clusters, err := GetAllAuroraFromStack(awsCtx, stackName)
	if err != nil {
		return "", err
	}

	// 対話的選択機能を使用
	selectedIndex, err := SelectFromOptions("複数のAuroraクラスターが見つかりました", clusters)
	if err != nil {
		return "", err
	}

	return clusters[selectedIndex], nil
}

// GetAllAuroraFromStack はスタック内のすべてのAuroraクラスター識別子を取得します
func GetAllAuroraFromStack(awsCtx AwsContext, stackName string) ([]string, error) {
	var results []string

	stackResources, err := getStackResources(awsCtx, stackName)
	if err != nil {
		return results, fmt.Errorf("CloudFormationスタックのリソース取得に失敗: %w", err)
	}

	// リソースの中からAurora DBClusterを探す
	for _, resource := range stackResources {
		if resource.ResourceType != nil && *resource.ResourceType == "AWS::RDS::DBCluster" {
			if resource.PhysicalResourceId != nil && *resource.PhysicalResourceId != "" {
				results = append(results, *resource.PhysicalResourceId)
				fmt.Printf("🔍 検出されたAuroraクラスター: %s\n", *resource.PhysicalResourceId)
			}
		}
	}

	if len(results) == 0 {
		return results, fmt.Errorf("指定されたスタック (%s) にAuroraクラスターが見つかりませんでした", stackName)
	}

	return results, nil
}
