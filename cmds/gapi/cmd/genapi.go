/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/

package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"path/filepath"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/spf13/cobra"
)

// genapiCmd represents the genapi command
var genapiCmd = &cobra.Command{
	Use:   "genapi resource[, resource]",
	Short: "生成API相关代码",
	Long: `生成包括路由、Handlers、业务模型等相关代码。例如：
	gapi genapi contents
会生成contents相关的：
	routers/contents.go
	internal/contents/handlers.go
	internal/contents/bizs.go
	......
等代码`,
	Args: func(cmd *cobra.Command, args []string) error {
		// 必须指定至少一个resource
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, resources []string) {
		// slog.Info("flags", "version", *genapiVersion, "bare", *genapiBare)
		// get module info
		modFile, err := modFileByFile()
		if err != nil {
			slog.Error("获取module信息失败", "error", err)
			return
		}
		mod := modInfo{
			Path: modFile.Module.Mod.Path,
		}

		// astTest("")
		// slog.Info("ast test")
		// return

		// range all resources
		for _, resource := range resources {
			rInfo := resourceInfo{
				Mod:               mod,
				Version:           *genapiVersion,
				Resource:          resource,
				Schema:            schemaInfo{Name: resourceSchemaName(resource)},
				ResourceBody:      resourceBody(resource),
				ResourcePatchBody: resourcePatchBody(resource),
			}

			// generate structure
			resourceFiles := []structureFile{}

			// codeTmpls
			codeTmpls := []codeTmpl{}

			// --bare
			if *genapiBare {
				codeTmpls = []codeTmpl{
					{tmpls.Routers_resources_bare, filepath.Join(ROUTER_BASE, resource+GO_EXT), rInfo},
				}
			}

			// --bare false (defalt)
			if !*genapiBare {
				resourceFiles = []structureFile{
					{true, filepath.Join(INTERNAL_BASE, resource), ""},
				}

				codeTmpls = []codeTmpl{
					// bizs
					{tmpls.Resource_bizs, filepath.Join(INTERNAL_BASE, resource, BIZS_BASE), rInfo},
					// handlers
					{tmpls.Resource_handlers, filepath.Join(INTERNAL_BASE, resource, HANDLERS_BASE), rInfo},
					// routers
					{tmpls.Resource_routers, filepath.Join(ROUTER_BASE, resource+GO_EXT), rInfo},
				}
			}

			// generate structure
			slog.Info("生成资源相关结构")
			if err := genStructure(resourceFiles); err != nil {
				slog.Error("生成资源相关结构失败", "error", err)
				return
			}

			// generate codes
			slog.Info("生成资源相关代码")
			if err := genCodes(codeTmpls); err != nil {
				slog.Error("生成资源相关代码失败", "error", err)
				return
			}
		}

	},
}

var astTest = parseSchema

func parseSchema(resource string) (*string, error) {
	// generate test

	// parse schema file
	fset := token.NewFileSet()
	filename := filepath.Join(INTERNAL_BASE, SCHEMAS_BASE, resource+GO_EXT)
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// find schema struct
	var schemaStruct *ast.StructType
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			schemaStruct = x.Type.(*ast.StructType)
			return false
		}
		return true
	})
	if schemaStruct == nil {
		return nil, errors.New("schema struct not found")
	}

	// range field list to find our field
	for _, field := range schemaStruct.Fields.List {
		fmt.Println(field.Names, field.Type, field.Tag, field.Comment)
		switch t := field.Type.(type) {
		case *ast.Ident:
			fmt.Println(t.Name)
		case *ast.SelectorExpr:
			fmt.Println(t.X.(*ast.Ident).Name, t.Sel.Name)
		}
	}

	return nil, nil

}

func resourceBody(resource string) string {
	// parse resource schema

	return ""
}

func resourcePatchBody(resource string) string {
	return ""
}

var (
	genapiVersion *string
	genapiBare    *bool
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
	genapiBare = genapiCmd.Flags().BoolP("bare", "b", false, "是否为纯路由")
}
