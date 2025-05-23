package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Ec2Instance EC2インスタンスの情報を格納する構造体
type Ec2Instance struct {
	InstanceId   string
	InstanceName string
	State        string
}

// ListEc2Instances 現在のリージョンのEC2インスタンス一覧を取得する
func ListEc2Instances(awsCtx AwsContext) ([]Ec2Instance, error) {
	cfg, err := LoadAwsConfig(awsCtx)
	if err != nil {
		return nil, fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := ec2.NewFromConfig(cfg)
	result, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("EC2インスタンス一覧の取得に失敗: %w", err)
	}

	var instances []Ec2Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// 終了済みのインスタンスは除外
			if instance.State.Name == types.InstanceStateNameTerminated {
				continue
			}

			// インスタンス名を取得（Nameタグから）
			instanceName := "（名前なし）"
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" && tag.Value != nil {
					instanceName = *tag.Value
					break
				}
			}

			instances = append(instances, Ec2Instance{
				InstanceId:   *instance.InstanceId,
				InstanceName: instanceName,
				State:        string(instance.State.Name),
			})
		}
	}

	return instances, nil
}

// StartEc2Instance EC2インスタンスを起動する
func StartEc2Instance(awsCtx AwsContext, instanceId string) error {
	cfg, err := LoadAwsConfig(awsCtx)
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := ec2.NewFromConfig(cfg)
	_, err = client.StartInstances(context.Background(), &ec2.StartInstancesInput{
		InstanceIds: []string{instanceId},
	})
	if err != nil {
		return fmt.Errorf("❌ EC2インスタンスの起動に失敗: %w", err)
	}
	return nil
}

// StopEc2Instance EC2インスタンスを停止する
func StopEc2Instance(awsCtx AwsContext, instanceId string) error {
	cfg, err := LoadAwsConfig(awsCtx)
	if err != nil {
		return fmt.Errorf("AWS設定のロードに失敗: %w", err)
	}

	client := ec2.NewFromConfig(cfg)
	_, err = client.StopInstances(context.Background(), &ec2.StopInstancesInput{
		InstanceIds: []string{instanceId},
	})
	if err != nil {
		return fmt.Errorf("❌ EC2インスタンスの停止に失敗: %w", err)
	}
	return nil
}

// SelectInstanceInteractively EC2インスタンス一覧を表示してユーザーに選択させる
func SelectInstanceInteractively(awsCtx AwsContext) (string, error) {
	fmt.Println("EC2インスタンス一覧を取得中...")

	instances, err := ListEc2Instances(awsCtx)
	if err != nil {
		return "", fmt.Errorf("❌ EC2インスタンス一覧の取得に失敗: %w", err)
	}

	if len(instances) == 0 {
		return "", fmt.Errorf("❌ 利用可能なEC2インスタンスが見つかりません")
	}

	// インスタンス一覧を表示
	fmt.Println("\n📋 利用可能なEC2インスタンス:")
	fmt.Println("番号 | インスタンスID        | インスタンス名                | 状態")
	fmt.Println("-----|----------------------|------------------------------|----------")

	for i, instance := range instances {
		fmt.Printf("%-4d | %-20s | %-28s | %s\n",
			i+1, instance.InstanceId, instance.InstanceName, instance.State)
	}

	// ユーザーに選択させる
	fmt.Print("\n操作するインスタンスの番号を入力してください: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("❌ 入力の読み取りに失敗: %w", err)
	}

	// 入力を数値に変換
	input = strings.TrimSpace(input)
	selectedNum, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("❌ 無効な番号です: %s", input)
	}

	// 範囲チェック
	if selectedNum < 1 || selectedNum > len(instances) {
		return "", fmt.Errorf("❌ 番号は1から%dの間で入力してください", len(instances))
	}

	selectedInstance := instances[selectedNum-1]
	fmt.Printf("✅ 選択されたインスタンス: %s (%s)\n",
		selectedInstance.InstanceName, selectedInstance.InstanceId)

	return selectedInstance.InstanceId, nil
}
