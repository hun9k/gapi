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
		// mod info
		wd, err := os.Getwd()
		if err != nil {
			slog.Warn("get workdir error", "error", err)
			wd = "."
		}
		basePath := path.Base(args[0])
		mod := modInfo{
			Path: args[0],
			Base: basePath,
			Dir:  filepath.Join(wd, basePath),
		}

		// generate basic structure
		dirs := []string{
			mod.Dir,
			filepath.Join(mod.Dir, "apis"),
			filepath.Join(mod.Dir, "routers"),
			filepath.Join(mod.Dir, "models"),
			filepath.Join(mod.Dir, "handlers"),
		}
		slog.Info("开始创建基础目录")
		if err := genDirs(dirs); err != nil {
			slog.Error("创建基础目录失败", "error", err)
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
			{"# {{.Base}}\n", filepath.Join(mod.Dir, "README.md"), mod},
			{tmpls.Configs, filepath.Join(mod.Dir, "configs.yaml"), nil},
			{tmpls.RoutersPing, filepath.Join(mod.Dir, "routers", "ping.go"), mod},
			{tmpls.Main, filepath.Join(mod.Dir, "main.go"), mod},
		}
		slog.Info("开始生成基础代码")
		if err := genCodes(codeTmpls); err != nil {
			slog.Error("生成基础代码失败", "error", err)
			return
		}

		// success and tips
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
