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
	"strings"
	"time"

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
		// 不能在已初始化项目下执行
		if _, err := os.Stat(LOCK_FILE); err == nil {
			slog.Error("不能在已初始化项目下执行")
			return
		}

		// mod 信息
		modPath := args[0]
		modBase := path.Base(modPath)

		// 当前工作目录
		curPath, err := os.Getwd()
		if err != nil {
			slog.Error("工作目录失败", "error", err)
			return
		}
		appPath := filepath.Join(curPath, modBase)

		// 坚持模块目录是否存在
		if ex, err := DirExists(appPath); err != nil || ex {
			slog.Error("模块目录存在", "filename", appPath, "error", err)
			return
		}

		// 创建module目录
		if err := MkDir(filepath.Join(appPath)); err != nil {
			slog.Error("创建目录失败", "filename", appPath, "error", err)
			return
		}
		slog.Info("创建目录成功", "filename", appPath)

		// change dir to mod dir
		if err := os.Chdir(appPath); err != nil {
			slog.Error("切换目录失败", "error", err)
			return
		}
		slog.Info("切换目录成功", "current", appPath)

		// mod init
		cmdInit := exec.Command("go", "mod", "init", modPath)
		if err := cmdInit.Run(); err != nil {
			slog.Error("模块初始失败", "cmd", "go mod init "+modPath, "error", err)
			return
		}
		if err := exec.Command("go", "mod", "edit", `-replace=github.com/hun9k/gapi=D:/apps/gapi`).Run(); err != nil {
			slog.Error("模块初始失败", "cmd", "", "error", err)
			return
		}
		slog.Info("模块初始成功", "cmd", "go mod init "+modPath)

		// 获取mod信息filename
		mod, err := ModFile(filepath.Join(appPath, "go.mod"))
		if err != nil {
			slog.Error("获取模块失败", "error", err)
			return
		}
		slog.Info("获取模块成功", "filename", filepath.Join(appPath, "go.mod"))

		// 创建基础目录
		dirs := []string{
			filepath.Join(appPath, HANDLER_DIR),
			filepath.Join(appPath, TASK_DIR),
			filepath.Join(appPath, MODEL_DIR),
		}
		for i, err := range MkDirs(dirs) {
			if err != nil {
				slog.Error("创建目录失败", "filename", dirs[i], "error", err)
				continue
			}
			slog.Info("创建目录成功", "filename", dirs[i])
		}

		// 生成基础代码
		codeTmpls := []fileTmpl{
			{
				text:     "{{.initTime}}",
				filename: filepath.Join(appPath, LOCK_FILE),
				data: tmplData{
					"initTime": time.Now().Format(time.RFC3339),
				},
			},
			{
				text:     "# {{.modBase}}\n",
				filename: filepath.Join(appPath, README_FILE),
				data: tmplData{
					"modBase": modBase,
				},
			},
			{
				text:     tmpls.Dockerfile,
				filename: filepath.Join(appPath, DOCKERFILE_FILE),
			},
			{
				text:     tmpls.DockerCompose,
				filename: filepath.Join(appPath, DOCKERCOMPOSE_FILE),
				data: tmplData{
					"dbName": modBase,
				},
			},
			{
				text:     tmpls.GitIgnore,
				filename: filepath.Join(appPath, GITIGNORE_FILE),
			},
			{
				text:     tmpls.Configs,
				filename: filepath.Join(appPath, CONFIG_FILE),
				data: tmplData{
					"appName": modBase,
				},
			},
			{
				text:     tmpls.HandlersInit,
				filename: filepath.Join(appPath, HANDLER_DIR, INIT_FILE),
			},
			{
				text:     tmpls.Main,
				filename: filepath.Join(appPath, MAIN_FILE),
				data: tmplData{
					"modPath": mod.Module.Mod.Path,
				},
			},
		}
		for i, err := range GenFiles(codeTmpls) {
			if err != nil {
				slog.Error("生成文件失败", "filename", codeTmpls[i].filename, "error", err)
				continue
			}
			slog.Info("生成文件成功", "filename", codeTmpls[i].filename)
		}

		// 生成用户模块
		if strings.ToLower(*initUser) != "off" {
			// user model and handler
			if err := GenUserModelHandler(appPath, mod, *initUser); err != nil {
				return
			}
		}

		// mod tidy
		if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
			slog.Error("模块整理失败", "cmd", "go mod tidy", "error", err)
			return
		}
		slog.Info("模块整理成功", "cmd", "go mod tidy")
		slog.Info("初始化完成!!!")
	},
}

var (
	initUser *string
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initUser = initCmd.Flags().StringP("user", "u", "", "用户模块路径，off表示不启用user模块, 默认空表示不区分平台，其他例如admin表示admin平台下")
}
