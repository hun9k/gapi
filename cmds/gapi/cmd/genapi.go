/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/

package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
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
		// get module info
		modFile, err := modFileByFile()
		if err != nil {
			slog.Error("get mod info failed", "error", err)
			return
		}
		mod := modInfo{
			Path: modFile.Module.Mod.Path,
		}

		// range all resources
		resources := args
		for _, resource := range resources {
			rInfo := resourceInfo{
				Mod:      mod,
				Version:  *genapiVersion,
				Resource: resource,
			}

			// generate structure
			resourceFiles := []structureFile{}

			// codeTmpls
			codeTmpls := []codeTmpl{}

			// --bare
			if *genapiBare {
				codeTmpls = []codeTmpl{
					{tmpls.Routers_resources_bare, filepath.Join(ROUTER_BASE, resource+".go"), rInfo},
				}
			}

			// --bare false (defalt)
			if !*genapiBare {
				resourceFiles = []structureFile{
					{true, filepath.Join(INTERNAL_BASE, resource), ""},
				}

				codeTmpls = []codeTmpl{
					{tmpls.Routers_resources, filepath.Join(ROUTER_BASE, resource+".go"), rInfo},
					{tmpls.Resource_bizs, filepath.Join(INTERNAL_BASE, resource, "bizs.go"), rInfo},
					{tmpls.Resource_handlers, filepath.Join(INTERNAL_BASE, resource, "handlers.go"), rInfo},
				}
			}

			// generate structure
			if err := genStructure(resourceFiles); err != nil {
				slog.Error("generate resource structure failed", "error", err)
				return
			}

			// generate codes
			if err := genCodes(codeTmpls); err != nil {
				slog.Error("generate resource codes failed", "error", err)
				return
			}
		}

	},
}

func modFileByFile() (*modfile.File, error) {
	goModFilename := "go.mod"
	// 读取 go.mod 文件内容
	modBytes, err := os.ReadFile(goModFilename)
	if err != nil {
		return nil, err
	}

	// 解析 go.mod 内容
	modFile, err := modfile.Parse(goModFilename, modBytes, nil)
	if err != nil {
		return nil, err
	}

	return modFile, nil
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
