/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genapiCmd represents the genapi command
var genapiCmd = &cobra.Command{
	Use:   "genapi resource[, resource]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples and 
usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// 必须指定至少一个resource
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genapi called", args, cmd.Flags)
	},
}

var (
	genapiVersion *string
	genapiBare    *bool
	genapiCrud    *bool
)

func init() {
	rootCmd.AddCommand(genapiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	genapiVersion = genapiCmd.Flags().StringP("version", "v", "v1", "路由的版本号前缀")
	genapiBare = genapiCmd.Flags().BoolP("bare", "b", false, "是否纯路由")
	genapiCrud = genapiCmd.Flags().BoolP("crud", "d", false, "是否包含CRUD")
}
