/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/hun9k/gapi/conf"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init module-path",
	Short: "初始化module",
	Long: `用于初始化新module，会生成module的目录结构和基础代码。示例：
	gapi init github.com/hun9k/gapi-demo
会在gapi-demo目录创建module，module-path为github.com/hun9k/gapi-demo`,
	// Args: cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
		// 必须指定一个module-path
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		// check module-path
		// if !utils.IsValidModulePath(args[0]) {
		// 	return errors.New("module-path is invalid")
		// }
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// mod info
		mod := modInfo{
			Path: args[0],
			Base: path.Base(args[0]),
		}

		// generate basic structure
		basicFiles := []structureFile{
			{true, mod.Base, ""},
			{false, filepath.Join(mod.Base, README_BASE), fmt.Sprintf("# %s\n", mod.Path)},
			{true, filepath.Join(mod.Base, ROUTER_BASE), ""},
			{false, filepath.Join(mod.Base, ROUTER_BASE, README_BASE), "# 自定义路由在此\n"},
			{true, filepath.Join(mod.Base, INTERNAL_BASE), ""},
			{true, filepath.Join(mod.Base, INTERNAL_BASE, SCHEMAS_BASE), ""},
			{false, filepath.Join(mod.Base, INTERNAL_BASE, SCHEMAS_BASE, README_BASE), "# 表结构设计在此\n"},
		}
		slog.Info("开始生成基础结构")
		if err := genStructure(basicFiles); err != nil {
			slog.Error("生成基础结构失败", "error", err)
			return
		}

		// change dir to mod dir
		slog.Info("开始切换目录")
		if err := changeDir(mod); err != nil {
			slog.Error("切换目录失败", "error", err)
			return
		}

		// mod init
		slog.Info("开始go mod init")
		if err := modInit(mod); err != nil {
			slog.Error("go mod init执行失败", "error", err)
			return
		}

		// generate basic codes
		codeTmpls := []codeTmpl{
			{tmpls.Configs, CONFIGS_BASE, conf.Default()},
			{tmpls.Apis_group, filepath.Join(ROUTER_BASE, GROUP_ROUTER_BASE), mod},
			{tmpls.Main, MAIN_BASE, mod},
		}
		slog.Info("开始生成基础代码")
		if err := genCodes(codeTmpls); err != nil {
			slog.Error("生成基础代码失败", "error", err)
			return
		}

		// mod tidy
		// slog.Info("开始执行go mod tidy")
		// if err := modTidy(mod); err != nil {
		// 	slog.Error("go mod tidy执行失败", "error", err)
		// 	return
		// }

		// success and tips
		slog.Info("新module初始化于当前目录", "module-path", mod.Path)
		slog.Info("通常需要修改配置文件以适应环境", "file", "configs.yaml")
		slog.Info("接下来应执行", "cmd", fmt.Sprintf("cd %s && go mod tidy && go run .", mod.Base))
	},
}

// create module dir and change work dir to it
func changeDir(mod modInfo) error {
	// change work dir to module path
	if err := os.Chdir(mod.Base); err != nil {
		return err
	}

	return nil
}

// mod init
func modInit(mod modInfo) error {
	// execute `go mod init`
	cmdModInit := exec.Command("go", "mod", "init", mod.Path)
	if err := cmdModInit.Run(); err != nil {
		return err
	}

	// for dev, defer delete this command
	//go mod edit -replace="github.com/hun9k/gapi"="D:/apps/gapi"
	cmd1 := exec.Command("go", "mod", "edit", `-replace=github.com/hun9k/gapi=D:/apps/gapi`)
	if err := cmd1.Run(); err != nil {
		return err
	}
	return nil
}

// mod tidy
func modTidy(mod modInfo) error {
	// execute `go mod tidy`
	// cmd := exec.Command("go", "mod", "tidy")
	// if err := cmd.Run(); err != nil {
	// 	return err
	// }

	// cmd2 := exec.Command("go", "get", "github.com/hun9k/gapi/cmds/gapi/cmd@v0.0.0-00010101000000-000000000000")
	// if err := cmd2.Run(); err != nil {
	// 	return err
	// }

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
