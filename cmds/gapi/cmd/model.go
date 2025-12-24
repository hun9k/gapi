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
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// modelCmd represents the genmodel command
var modelCmd = &cobra.Command{
	Use:   "model <resource>[, <resource>...]",
	Short: "生成资源model",
	Long: `生成一个或多个模型基本代码，包括模型类型，对应的钩子，模型的CRUD代码，等。例如：

model user, post

会生成models/user_gen.go, models/user.go包含对应的代码`,
	// 必须指定至少一个model-name
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否在项目根目录下执行
		if _, err := os.Stat(LOCK_FILE); err != nil {
			slog.Error("仅支持在gapi项目根目录下执行")
			return
		}

		appPath, err := os.Getwd()
		if err != nil {
			slog.Error("获取当前工作目录失败", "error", err)
			return
		}

		// 创建模型目录
		modelPath := filepath.Join(appPath, MODEL_DIR)
		if err := MkDir(modelPath); err != nil {
			slog.Error("创建模型目录：失败", "dir", MODEL_DIR, "error", err)
			return
		}
		slog.Info("创建模型目录", "dir", modelPath)

		// 生成模型文件
		for _, res := range args {
			GenModel(res)
		}

		// 更新模型init
		SetModelInit(modelPath)

		// mod tidy
		if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
			slog.Error("mod tidy失败", "error", err)
			return
		}
	},
}

func getModels(path string) ([]string, error) {
	var result []string
	// 读取目录内容
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	// 遍历目录中的每个条目
	for _, entry := range entries {
		// 只处理文件，跳过目录
		if entry.IsDir() {
			continue
		}
		// 获取文件名
		filename := entry.Name()
		// 跳过不以 ".go" 结尾的文件
		if filepath.Ext(filename) != ".go" {
			continue
		}
		// 跳过以 "_init" 开头的文件
		if strings.HasPrefix(filename, "init") {
			continue
		}

		// 去掉文件后缀
		nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

		// 添加到结果中
		result = append(result, fmt.Sprintf("&%s{}", strcase.ToCamel(nameWithoutExt)))
	}

	return result, nil
}

func init() {
	rootCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genmodelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
