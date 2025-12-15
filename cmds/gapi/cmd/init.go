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
		// mod info
		modPath := args[0]
		modBase := path.Base(modPath)

		// check if module dir exists
		if ex, err := dirExists(modBase); err != nil || ex {
			slog.Error("模块目录已存在", "module", modPath)
			return
		}

		// 创建基础结构
		dirs := []string{
			filepath.Join(modBase),
			filepath.Join(modBase, "models"),
			filepath.Join(modBase, "handlers"),
		}
		for i, err := range genDirs(dirs) {
			if err != nil {
				slog.Error("创建目录失败", "dir", dirs[i], "error", err)
				continue
			}
			slog.Info("创建目录成功", "dir", dirs[i])
		}

		// 生成基础代码
		codeTmpls := []codeTmpl{
			{"# {{.modePath}}\n", filepath.Join(modBase, "README.md"), tmplData{"modePath": modPath}},
			{tmpls.Routers, filepath.Join(modBase, "handlers", "init_.go"), nil},
			{tmpls.Configs, filepath.Join(modBase, "configs.yaml"), nil},
			{tmpls.Main, filepath.Join(modBase, "main.go"), tmplData{"modPath": modPath}},
			{"", filepath.Join(modBase, ".gapi.lock"), nil},
		}
		for i, err := range genCodes(codeTmpls, true) {
			if err != nil {
				slog.Error("生成代码失败", "filename", codeTmpls[i].filename, "error", err)
				continue
			}
			slog.Info("生成代码成功", "filename", codeTmpls[i].filename)
		}

		// change dir to mod dir
		if err := changeDir(modBase); err != nil {
			slog.Error("切换模块目录失败", "error", err)
			return
		}

		// mod init
		if err := modInit(modPath); err != nil {
			slog.Error("mod init失败", "error", err)
			return
		}
		slog.Info("mod init成功")

		if err := modTidy(); err != nil {
			slog.Error("mod tidy失败", "error", err)
			return
		}
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

// mod init and tidy
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

func modTidy() error {
	cmdModTidy := exec.Command("go", "mod", "tidy")
	if err := cmdModTidy.Run(); err != nil {
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
