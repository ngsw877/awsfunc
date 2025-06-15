package cmd

import (
	"awstk/internal/aws"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

const AppName = "awstk"

var region string
var profile string
var awsCtx aws.Context
var stackName string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/." + AppName + ".yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVarP(&region, "region", "R", "ap-northeast-1", "AWSリージョン")
	RootCmd.PersistentFlags().StringVarP(&profile, "profile", "P", "", "AWSプロファイル")
	RootCmd.PersistentFlags().StringVarP(&stackName, "stack", "S", "", "CloudFormationスタック名")

	// コマンド実行前に共通でプロファイルチェックとawsCtx設定を行う
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return nil
		}

		// プロファイルチェック
		err := checkAndSetProfile(cmd)
		if err != nil {
			return err
		}

		// awsCtxを設定
		awsCtx = aws.Context{Region: region, Profile: profile}

		return nil
	}
}

// checkAndSetProfile はプロファイルの確認と設定を行うプライベート関数
func checkAndSetProfile(cmd *cobra.Command) error {
	// プロファイルがすでに指定されている場合は何もしない
	if profile != "" {
		return nil
	}
	// 環境変数からプロファイル取得を試みる
	envProfile := os.Getenv("AWS_PROFILE")
	if envProfile == "" {
		// プロファイルが見つからない場合はエラー
		cmd.SilenceUsage = true // エラー時のUsage表示を抑制
		return errors.New("❌ エラー: プロファイルが指定されていません。-Pオプションまたは AWS_PROFILE 環境変数を指定してください")
	}
	// 環境変数からプロファイルを設定
	profile = envProfile
	cmd.Println("🔍 環境変数 AWS_PROFILE の値 '" + profile + "' を使用します")
	return nil
}
