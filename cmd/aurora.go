package cmd

import (
	"awsfunc/internal"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	auroraClusterId string
	auroraProfile   string
)

var auroraCmd = &cobra.Command{
	Use:   "aurora",
	Short: "Aurora関連の操作を行うコマンド群",
	Long:  "AWS Aurora DBクラスターの操作を行うCLIコマンド群です。",
}

var auroraStartClusterCmd = &cobra.Command{
	Use:   "start",
	Short: "Aurora DBクラスターを起動する",
	Long: `指定したAurora DBクラスターを起動します。

例:
  awsfunc aurora start-cluster -d <aurora-cluster-identifier> [-P <aws-profile>]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if auroraClusterId == "" {
			return fmt.Errorf("❌ Aurora DBクラスター識別子は必須です")
		}
		fmt.Printf("Aurora DBクラスター (%s) を起動します...\n", auroraClusterId)

		err := internal.StartAuroraCluster(auroraClusterId, auroraProfile)

		if err != nil {
			fmt.Printf("❌ Aurora DBクラスターの起動に失敗しました。")
		}

		fmt.Println("✅ Aurora DBクラスターの起動を開始しました。起動完了まで数十分かかります。")
		return nil
	},
	SilenceUsage: true,
}

var auroraStopClusterCmd = &cobra.Command{
	Use:   "stop",
	Short: "Aurora DBクラスターを停止する",
	Long: `指定したAurora DBクラスターを停止します。

例:
  awsfunc aurora stop-cluster -d <aurora-cluster-identifier> [-P <aws-profile>]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if auroraClusterId == "" {
			return fmt.Errorf("❌ Aurora DBクラスター識別子は必須です")
		}
		fmt.Printf("Aurora DBクラスター (%s) を停止します...\n", auroraClusterId)

		err := internal.StopAuroraCluster(auroraClusterId, auroraProfile)
		if err != nil {
			fmt.Printf("❌ Aurora DBクラスターの停止に失敗しました。")
			return err
		}

		fmt.Println("✅ Aurora DBクラスターの停止を開始しました。")
		return nil
	},
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(auroraCmd)
	auroraCmd.AddCommand(auroraStartClusterCmd)
	auroraCmd.AddCommand(auroraStopClusterCmd)
	auroraCmd.PersistentFlags().StringVarP(&auroraClusterId, "db-cluster-identifier", "d", "", "Aurora DBクラスター識別子（必須）")
}
