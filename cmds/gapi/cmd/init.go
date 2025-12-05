/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <module-path>",
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

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// curr work dir
		wd, err := os.Getwd()
		if err != nil {
			slog.Warn("get workdir error", "error", err)
			wd = "."
		}

		// mod info
		modPath := args[0]
		modDir := filepath.Join(wd, path.Base(modPath))

		// mk module dir
		slog.Info("创建模块目录")
		moduleDir := []string{
			modDir,
		}
		if err := genDirs(moduleDir); err != nil {
			slog.Error("创建模块目录失败", "error", err)
			return
		}
		// change dir to mod dir
		slog.Info("切换模块目录")
		if err := changeDir(modDir); err != nil {
			slog.Error("切换模块目录失败", "error", err)
			return
		}

		// mod init
		slog.Info("初始化模块")
		if err := modInit(modPath); err != nil {
			slog.Error("初始化模块失败", "error", err)
			return
		}

		// make module dir
		dirs := []string{
			filepath.Join(modDir, "apis"),
			filepath.Join(modDir, "routers"),
			filepath.Join(modDir, "models"),
			filepath.Join(modDir, "handlers"),
		}

		slog.Info("创建基本结构")
		if err := genDirs(dirs); err != nil {
			slog.Error("创建基本结构失败", "error", err)
			return
		}

		// generate basic codes
		codeTmpls := []codeTmpl{
			{"# {{.Base}}\n", filepath.Join(modDir, "README.md"), tmplData{"Base": modPath}},
			{tmpls.Configs, filepath.Join(modDir, "configs.yaml"), nil},
			{tmpls.RoutersPing, filepath.Join(modDir, "routers", "ping.go"), nil},
			{tmpls.Main, filepath.Join(modDir, "main.go"), tmplData{"Path": modPath}},
		}
		slog.Info("生成基础代码")
		if err := genCodes(codeTmpls); err != nil {
			slog.Error("生成基础代码失败", "error", err)
			return
		}

		// success and tips
		slog.Info("模块初始化成功，可以初始化代码版本库了")
		slog.Info("接下来，创建openapi文件，生成代码")

	},
}

// create module dir and change work dir to it
func changeDir(dir string) error {
	// change work dir to module path
	if err := os.Chdir(dir); err != nil {
		return err
	}

	return nil
}

// mod init
func modInit(path string) error {
	// execute `go mod init`
	cmdModInit := exec.Command("go", "mod", "init", path)
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
