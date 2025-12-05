/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

// genhandlerCmd represents the genhandler command
var genhandlerCmd = &cobra.Command{
	Use:   "genhandler",
	Short: "生成Handler代码",
	Long: `基于openapi3.x生成API的Router和Handler。示例：

genhandler -f ./openapi.yaml -o ./handler

使用openapi.yaml定义的API文件，生成的Router和Handler代码相关代码。`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("genhandler called")

		// 解析openapi3.x文件
		doc, err := parseOpenapiFile(*genhandlerFilename)
		if err != nil {
			slog.Error("parseOpenapiFile error", slog.Any("err", err))
			return
		}

		// 生成schemas
		if err := genSchemas(doc, "internal/messages/schemas_gen.go"); err != nil {
			slog.Error("genSchemas error", slog.Any("err", err))
			return
		}

		if err := genPathMessages(doc, "/pet", "internal/messages"); err != nil {
			slog.Error("genMessages error", slog.Any("err", err))
			return
		}

		// 生成handlers
		if err := genHandlers(doc, "internal/routers/routers_gen.go"); err != nil {
			slog.Error("genHandler error", slog.Any("err", err))
			return
		}

	},
}

const (
	P_GAPI_HTTP = "github.com/hun9k/gapi/http"
	P_GIN       = "github.com/gin-gonic/gin"
)

func parseOpenapiFile(filename string) (*openapi3.T, error) {
	// 创建解析器
	loader := openapi3.NewLoader()
	// 可设置严格模式（验证规范合法性）
	loader.IsExternalRefsAllowed = true
	// 解析为OpenAPI文档对象
	doc, err := loader.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}

	// 严格校验版本：仅支持 3.0.x（3.0.0/3.0.1/3.0.2/3.0.3）
	if !strings.HasPrefix(doc.OpenAPI, "3.0.") {
		return nil, fmt.Errorf("3.0.x version only, current version is %s", doc.OpenAPI)
	}

	// 验证文档的规范合法性
	if err := doc.Validate(loader.Context); err != nil {
		return nil, err
	}

	return doc, nil
}

var (
	genhandlerFilename, genhandlerOutput *string
)

func init() {
	rootCmd.AddCommand(genhandlerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genhandlerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genhandlerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genhandlerFilename = genhandlerCmd.Flags().StringP("filename", "f", "./apis/openapi.yaml", "openapi3文件")
	genhandlerOutput = genhandlerCmd.Flags().StringP("output", "o", "./handlers", "输出文件目录")
}
