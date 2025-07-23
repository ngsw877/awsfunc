package cmd

import (
	"awstk/internal/aws"
	"errors"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

const AppName = "awstk"

var region string
var profile string
var awsCfg awsconfig.Config
var stackName string
var cfnClient *cloudformation.Client
var rdsClient *rds.Client

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: "AWS リソース管理用 CLI ツール",
	Long: `awstk は AWS リソースを効率的に管理するための CLI ツールです。

S3、ECR、ECS、CloudFormation などの各種 AWS サービスに対して、
一括削除や状態確認などの便利な操作を提供します。

使用例:
  awstk cleanup all -k "test"    # "test"を含むS3/ECRを一括削除
  awstk s3 gunzip my-bucket/logs # S3の.gzファイルを一括ダウンロード&解凍
  awstk ecs exec -s my-service   # Fargateコンテナへシェル接続`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// isAuthNotRequired は認証が不要なコマンドかどうかを判定する
func isAuthNotRequired(cmd *cobra.Command) bool {
	// 認証が不要なコマンド
	if cmd.Name() == "help" ||
		cmd.Name() == "version" {
		return true
	}
	// 認証不要なコマンドのサブコマンド
	if cmd.Parent() != nil &&
		cmd.Parent().Name() == "env" {
		return true
	}
	return false
}

// checkProfile はプロファイルの確認のみを行うプライベート関数
func checkProfile(cmd *cobra.Command) error {
	// プロファイルがすでに指定されている場合は案内を出して終了
	if profile != "" {
		cmd.Println("🔍 -Pオプションで指定されたプロファイル '" + profile + "' を使用します")
		return nil
	}
	// 環境変数からプロファイル取得を試みる
	envProfile := os.Getenv("AWS_PROFILE")
	if envProfile == "" {
		// プロファイルが見つからない場合はエラー
		cmd.SilenceUsage = true // エラー時のUsage表示を抑制
		return errors.New("❌ エラー: プロファイルが指定されていません。-Pオプションまたは AWS_PROFILE 環境変数を指定してください")
	}
	cmd.Println("🔍 環境変数 AWS_PROFILE の値 '" + envProfile + "' を使用します")
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&region, "region", "R", "ap-northeast-1", "AWSリージョン")
	RootCmd.PersistentFlags().StringVarP(&profile, "profile", "P", "", "AWSプロファイル")

	// コマンド実行前に共通でプロファイルチェックとawsCtx設定を行う
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// 認証が不要なコマンドはスキップ
		if isAuthNotRequired(cmd) {
			return nil
		}

		// プロファイルチェック
		err := checkProfile(cmd)
		if err != nil {
			return err
		}

		// awsCtxを設定
		awsCtx := aws.Context{Region: region, Profile: profile}

		// AWS設定を読み込み
		awsCfg, err = aws.LoadAwsConfig(awsCtx)
		if err != nil {
			return fmt.Errorf("AWS設定の読み込みエラー: %w", err)
		}

		return nil
	}
}
