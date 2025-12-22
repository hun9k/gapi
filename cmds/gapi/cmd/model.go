/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
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

		// models dir
		dirs := []string{
			filepath.Join(MODEL_DIR),
		}
		for i, err := range genDirs(dirs) {
			if err != nil {
				slog.Error("创建目录失败", "dir", dirs[i], "error", err)
				continue
			}
			slog.Info("创建目录成功", "dir", dirs[i])
		}

		codeTmpls := make([]codeTmpl, len(args))
		for i, resource := range args {
			codeTmpls[i] = codeTmpl{
				text:     tmpls.ResourceModel,
				filename: filepath.Join(MODEL_DIR, resource+GO_EXT),
				data: tmplData{
					"package":  MODEL_DIR,
					"model":    strcase.ToCamel(resource),
					"resource": resource,
				},
			}
		}

		// generate codes
		for i, err := range genCodes(codeTmpls, *mf_force) {
			if err != nil {
				slog.Error("生成模型失败", "filename", codeTmpls[i].filename, "error", err)
				continue
			}
			slog.Info("生成模型成功", "filename", codeTmpls[i].filename)
		}

		// generate migrate init
		modelList, err := getModels(MODEL_DIR)
		if err != nil {
			slog.Error("获取模型名列表失败", "error", err)
			return
		}
		codeTmpls = []codeTmpl{
			{
				text:     tmpls.ModelsInit,
				filename: filepath.Join(MODEL_DIR, "init.go"),
				data: tmplData{
					"modelList": strings.Join(modelList, ", "),
				},
			},
		}
		for i, err := range genCodes(codeTmpls, true) {
			if err != nil {
				slog.Error("生成模型init失败", "filename", codeTmpls[i].filename, "error", err)
				continue
			}
			slog.Info("生成模型init成功", "filename", codeTmpls[i].filename)
		}

		// mod tidy
		if err := modTidy(); err != nil {
			slog.Error("mod tidy失败", "error", err)
			return
		}
	},
}

func getModels(dir string) ([]string, error) {
	var result []string
	// 读取目录内容
	entries, err := os.ReadDir(dir)
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

var (
	mf_force *bool
	mf_op    *string
)

func init() {
	rootCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genmodelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mf_force = modelCmd.Flags().BoolP("force", "f", false, "文件存在时也强制生成")
	mf_op = modelCmd.Flags().StringP("op", "o", "create", "操作类型，create，或init")
}
