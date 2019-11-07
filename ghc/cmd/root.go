package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ghc "github.com/youyo/github-actions-comment"
)

var Version string

var rootCmd = &cobra.Command{
	Use:          "ghc",
	Short:        "",
	Version:      Version,
	RunE:         ghc.Run,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("title", "t", "Result", "comment title")
	rootCmd.Flags().StringP("body", "b", "", "comment body")
	rootCmd.Flags().BoolP("failure", "f", false, "set failure status")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {}
