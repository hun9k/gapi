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
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// handlerCmd represents the genhandler command
var handlerCmd = &cobra.Command{
	Use:   "handler <resource>[, <resource>...]",
	Short: "基于资源handler",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// 至少指定一个resource
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否在项目根目录下执行
		if _, err := os.Stat(LOCK_FILE); err != nil {
			slog.Error("仅支持在gapi项目根目录下执行")
			return
		}

		// 检查是否指定了平台
		if *hf_plat == "" && !*hf_splat {
			slog.Error("请指定平台-p <platform>，或使用单平台-s")
			return
		}

		// 读取go.mod
		mod, err := modFile("go.mod")
		if err != nil {
			slog.Error("读取go.mod失败", "error", err)
			return
		}

		for _, resource := range args {
			// 生成目录
			dirs := []string{
				filepath.Join("handlers", *hf_plat, resource),
			}
			for i, err := range genDirs(dirs) {
				if err != nil {
					slog.Error("创建目录失败", "dir", dirs[i], "error", err)
					continue
				}
				slog.Info("创建目录成功", "dir", dirs[i])
			}

			// 解析模型结构

			// 生成messages
			// 生成handlers
			// 生成routers
			modelName := "models." + strcase.ToCamel(resource)
			codeTmpls := []codeTmpl{
				{tmpls.ResourceMessages, filepath.Join("handlers", *hf_plat, resource, "messages.go"), tmplData{"resource": resource}},
				{tmpls.ResourceHandlers, filepath.Join("handlers", *hf_plat, resource, "handlers.go"), tmplData{"resource": resource, "modelName": modelName, "modPath": mod.Module.Mod.Path}},
				{tmpls.ResourceRouters, filepath.Join("handlers", *hf_plat, resource, "routers.go"), tmplData{"resource": resource}},
			}
			for i, err := range genCodes(codeTmpls, *hf_force) {
				if err != nil {
					slog.Error("生成代码失败", "filename", codeTmpls[i].filename, "error", err)
					continue
				}
				slog.Info("生成代码成功", "filename", codeTmpls[i].filename)
			}
			// 注册路由
		}

		if err := modTidy(); err != nil {
			slog.Error("mod tidy失败", "error", err)
			return
		}
	},
}

// func setupRouter(resource string) error {
// 	return nil
// }

var (
	hf_plat  *string
	hf_splat *bool
	hf_force *bool
)

func init() {
	rootCmd.AddCommand(handlerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genhandlerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genhandlerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	hf_force = handlerCmd.Flags().BoolP("force", "f", false, "强制生成")
	hf_splat = handlerCmd.Flags().BoolP("splat", "s", false, "单平台")
	hf_plat = handlerCmd.Flags().StringP("platform", "p", "", "选择平台")
}
