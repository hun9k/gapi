/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newCmd represents the init command
var newCmd = &cobra.Command{
	Use:   "new module-path",
	Short: "New an APIs module",
	Long: `
gapi new github.com/hun9k/api-demo

A module will be generated in the api-demo directory, and module-path will be github.com/hun9k/api-demo.
Use --dir path to specify the other directory.
`,
	// Args: cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return errors.New("exactly one module-path must be specified")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called", args, cmd.Flag("dir").Value)

		// evalate module directory
		dir := cmd.Flag("dir").Value.String()
		if dir == "" {
			dir = path.Base(args[0])
		}

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		// mkdir
		moduleDir := filepath.Join(wd, dir)
		if err := os.Mkdir(moduleDir, 0750); err != nil {
			if os.IsExist(err) {
				fmt.Printf("%s already exists\n", moduleDir)
			} else {
				fmt.Println(err)
			}
			return
		}
		fmt.Printf("module directory %s has made.", moduleDir)

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().StringP("dir", "d", "", "specify the module directory, if not specified, the basename of the module-path will be used")
}
