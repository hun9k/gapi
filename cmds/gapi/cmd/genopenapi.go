/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// genopenapiCmd represents the genopenapi command
var genopenapiCmd = &cobra.Command{
	Use:   "genopenapi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := genOpenApi()
		if err != nil {
			slog.Error("生成OpenApi失败", "error", err)
			return
		}

		os.WriteFile("./apis/openapi.yaml", content, FILE_MODE)
	},
}

func init() {
	rootCmd.AddCommand(genopenapiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genopenapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genopenapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
