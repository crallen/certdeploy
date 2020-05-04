package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crallen/certdeploy/deploy"
	"github.com/spf13/cobra"
)

var (
	configFile string
	kubeConfig string
)

var rootCmd = &cobra.Command{
	Use:   "certdeploy",
	Short: "Deploy certificates to multiple Kubernetes clusters",
	RunE: func(cmd *cobra.Command, args []string) error {
		runner, err := deploy.New(configFile, kubeConfig)
		if err != nil {
			return err
		}
		return runner.Run()
	},
}

func init() {
	homeDir, _ := os.UserHomeDir()
	kubeConfigDefault := filepath.Join(homeDir, ".kube", "config")

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "k", kubeConfigDefault, "path to kubeconfig file")

	rootCmd.MarkFlagRequired("config")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
