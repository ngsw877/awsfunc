package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// StartAuroraCluster Auroraクラスタを起動する
func StartAuroraCluster(clusterID, profile string) error {
	ctx := context.Background()
	var cfg aws.Config
	var err error

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := rds.NewFromConfig(cfg)
	_, err = client.StartDBCluster(ctx, &rds.StartDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
	})
	if err != nil {
		return fmt.Errorf("Aurora DBクラスターの起動に失敗: %w", err)
	}
	return nil
}

// StopAuroraCluster Auroraクラスタを停止する
func StopAuroraCluster(clusterID, profile string) error {
	ctx := context.Background()
	var cfg aws.Config
	var err error

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := rds.NewFromConfig(cfg)
	_, err = client.StopDBCluster(ctx, &rds.StopDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
	})
	if err != nil {
		return fmt.Errorf("Aurora DBクラスターの停止に失敗: %w", err)
	}
	return nil
}
