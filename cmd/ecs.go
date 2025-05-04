package cmd

import (
	"awsfunc/internal"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	stackName     string
	clusterName   string
	serviceName   string
	containerName string
)

var EcsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "ECSリソース操作コマンド",
	Long:  `ECSリソースを操作するためのコマンド群です。`,
}

var ecsExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Fargateコンテナに接続するコマンド",
	Long: `Fargateコンテナにシェル接続するコマンドです。
CloudFormationスタック名を指定するか、クラスター名とサービス名を直接指定することができます。

例:
  awsfunc ecs exec -P my-profile -S my-stack
  awsfunc ecs exec -P my-profile -c my-cluster -s my-service -t app`,
	Run: func(cmd *cobra.Command, args []string) {
		// プロファイルチェック - 各internal関数内でLoadAWSConfigを実行するが、
		// 先にプロファイルが指定されているか確認する
		if Profile == "" && os.Getenv("AWS_PROFILE") != "" {
			Profile = os.Getenv("AWS_PROFILE")
			fmt.Println("🔍 環境変数 AWS_PROFILE の値 '" + Profile + "' を使用します")
		}

		if Profile == "" {
			fmt.Fprintln(os.Stderr, "❌ エラー: プロファイルが指定されていません。-Pオプションまたは AWS_PROFILE 環境変数を指定してください")
			os.Exit(1)
		}

		var cluster, service string

		// スタック名から情報取得
		if stackName != "" {
			fmt.Println("CloudFormationスタックからECS情報を取得します...")
			serviceInfo, err := internal.GetEcsFromStack(stackName, Region, Profile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
				os.Exit(1)
			}
			cluster = serviceInfo.ClusterName
			service = serviceInfo.ServiceName

			fmt.Println("🔍 検出されたクラスター: " + cluster)
			fmt.Println("🔍 検出されたサービス: " + service)
		} else if clusterName != "" && serviceName != "" {
			// クラスター名とサービス名が直接指定された場合
			cluster = clusterName
			service = serviceName
		} else {
			fmt.Fprintln(os.Stderr, "❌ エラー: スタック名が指定されていない場合は、クラスター名 (-c) とサービス名 (-s) が必須です")
			cmd.Help()
			os.Exit(1)
		}

		// タスクIDを取得
		taskId, err := internal.GetRunningTask(cluster, service, Region, Profile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
			os.Exit(1)
		}

		// シェル接続を実行
		fmt.Printf("🔍 コンテナ '%s' に接続しています...\n", containerName)
		err = internal.ExecuteCommand(cluster, taskId, containerName, Region, Profile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ コンテナへの接続に失敗しました: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(EcsCmd)
	EcsCmd.AddCommand(ecsExecCmd)

	// フラグを設定
	ecsExecCmd.Flags().StringVarP(&stackName, "stack", "S", "", "CloudFormationスタック名")
	ecsExecCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "ECSクラスター名 (-Sが指定されていない場合に必須)")
	ecsExecCmd.Flags().StringVarP(&serviceName, "service", "s", "", "ECSサービス名 (-Sが指定されていない場合に必須)")
	ecsExecCmd.Flags().StringVarP(&containerName, "container", "t", "app", "接続するコンテナ名")
}
