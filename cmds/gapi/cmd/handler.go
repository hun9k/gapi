/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

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

		// 读取go.mod
		mod, err := modFile("go.mod")
		if err != nil {
			slog.Error("读取go.mod失败", "error", err)
			return
		}

		// platform init
		// 生成目录
		dirs := []string{
			filepath.Join("handlers", *hf_plat),
		}
		for i, err := range genDirs(dirs) {
			if err != nil {
				slog.Error("创建平台目录失败", "dir", dirs[i], "error", err)
				continue
			}
			slog.Info("创建平台目录成功", "dir", dirs[i])
		}
		// 平台代码
		pSetup := "setup.go"
		if *hf_plat != "" {
			pSetup = *hf_plat + "_" + pSetup
		}
		pRouter := "routers_gen.go"
		if *hf_plat != "" {
			pRouter = *hf_plat + "_" + pRouter
		}
		codeTmpls := []codeTmpl{
			{
				text:     tmpls.PlatformSetup,
				filename: filepath.Join("handlers", pSetup),
				data:     tmplData{"platform": *hf_plat},
				isKeep:   true,
			},
			{
				text:     tmpls.PlatformRouters,
				filename: filepath.Join("handlers", pRouter),
				data:     tmplData{"platform": *hf_plat, "modPath": mod.Module.Mod.Path, "resources": getPlatRes(filepath.Join("handlers", *hf_plat))},
			},
		}
		for i, err := range genCodes(codeTmpls, *hf_force) {
			if err != nil {
				slog.Error("生成平台代码失败", "filename", codeTmpls[i].filename, "error", err)
				continue
			}
			slog.Info("生成平台代码成功", "filename", codeTmpls[i].filename)
		}

		for _, resource := range args {
			// 生成资源目录
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

			// 生成messages
			fields, err := ParseModel(resource)
			if err != nil {
				slog.Error("解析模型失败", "error", err)
				continue
			}
			// 生成handlers
			// 生成routers
			modelName := "models." + strcase.ToCamel(resource)
			codeTmpls := []codeTmpl{
				{
					text:     tmpls.ResourceMessages,
					filename: filepath.Join("handlers", *hf_plat, resource, "messages_gen.go"),
					data:     tmplData{"resource": resource, "fields": fields, "modelName": modelName, "modPath": mod.Module.Mod.Path, "iTime": iTime, "iSql": iSql, "iDatatypes": iDatatypes},
				},
				{
					text:     tmpls.ResourceHandlers,
					filename: filepath.Join("handlers", *hf_plat, resource, "handlers_gen.go"),
					data:     tmplData{"resource": resource, "modelName": modelName, "modPath": mod.Module.Mod.Path},
				},
				{
					text:     tmpls.ResourceRouters,
					filename: filepath.Join("handlers", *hf_plat, resource, "routers_gen.go"),
					data:     tmplData{"resource": resource},
				},
				{
					text:     tmpls.ResourceSetup,
					filename: filepath.Join("handlers", *hf_plat, resource, "setup.go"),
					data:     tmplData{"resource": resource},
					isKeep:   true,
				},
			}
			for i, err := range genCodes(codeTmpls, *hf_force) {
				if err != nil {
					slog.Error("生成代码失败", "filename", codeTmpls[i].filename, "error", err)
					continue
				}
				slog.Info("生成代码成功", "filename", codeTmpls[i].filename)
			}

			// 注册路由
			slog.Info("应在handlers/routers.go中添加路由")
		}

		if err := modTidy(); err != nil {
			slog.Error("mod tidy失败", "error", err)
			return
		}
	},
}

func getPlatRes(dir string) []string {
	var result []string
	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return result
	}
	// 遍历目录中的每个条目
	for _, entry := range entries {
		// 只处理文件，跳过目录
		if !entry.IsDir() || entry.Name() == "." || entry.Name() == ".." {
			continue
		}

		// 添加到结果中
		result = append(result, entry.Name())
	}

	return result
}

func ParseModel(res string) ([]mField, error) {
	filename := filepath.Join("models", res+".go")
	structName := strcase.ToCamel(res)

	// 2. 解析 Go 源文件，创建 AST
	fset := token.NewFileSet() // 文件集，用于位置信息
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fields := []mField{}

	// 3. 遍历 AST，查找目标结构体
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		// 我们只关心类型声明节点
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true // 继续遍历子节点
		}

		// 检查结构体名是否匹配
		if typeSpec.Name.Name != structName {
			return true // 名称不匹配，继续遍历
		}

		// 确认它是一个结构体类型
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false // 停止遍历此分支
		}

		// 找到了目标结构体！
		found = true

		// 4. 提取并打印字段信息
		fields = modelFields(structType)

		return false // 已找到，无需继续遍历
	})

	if !found {
		return nil, fmt.Errorf("未找到模型 %s", structName)
	}

	return fields, nil
}

type mField struct {
	Name string // 字段名
	Col  string // 列名
	Type string // 字段类型名
	Tag  string // 增加的Tag
	// Tag  string // 字段标签（原始字符串）
	// IsExported  bool   // 是否为导出字段（首字母大写）
	// IsAnonymous bool   // 是否为匿名字段
	IsAnonymous bool // 是否为匿名类型
	IsNonRef    bool // 是否为非引用类型
}

var iTime, iSql, iDatatypes bool

func modelFields(structType *ast.StructType) (fields []mField) {
	// non fileds
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		return []mField{}
	}
	// each field
	for _, field := range structType.Fields.List {
		// 匿名字段
		if field.Names == nil {
			// fields = append(fields, mField{
			// 	Type:        exprToString(field.Type),
			// 	IsAnonymous: true,
			// })
			continue
		}

		// 是否为引用
		isNonRef := true
		switch t := field.Type.(type) {
		case *ast.StarExpr, *ast.MapType:
			isNonRef = false
		case *ast.ArrayType:
			if _, err := parseArrayLength(t.Len); err != nil {
				isNonRef = false
			}
		}

		// title, body string
		for _, name := range field.Names {
			// tag中的列名
			col := strcase.ToSnake(name.Name)
			f := mField{
				Name:     name.Name,
				Col:      col,
				Type:     exprToString(field.Type),
				IsNonRef: isNonRef,
				Tag:      fmt.Sprintf("`json:\"%s\"`", col),
			}

			switch f.Type {
			case "time.Time", "*time.Time":
				iTime = true
			}

			fields = append(fields, f)
		}
	}

	return fields
}

// exprToString 是一个辅助函数，用于将 AST 表达式（类型）转换为字符串
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr: // 指针类型, e.g., *User
		return "*" + exprToString(t.X)
	case *ast.ArrayType: // 数组或切片类型, e.g., []string, [5]int
		l, err := parseArrayLength(t.Len)
		if err != nil {
			return "[]" + exprToString(t.Elt)
		} else {
			return fmt.Sprintf("[%d]%s", l, exprToString(t.Elt))
		}
	case *ast.MapType: // Map 类型, e.g., map[string]int
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.SelectorExpr: // 带包名的类型, e.g., time.Time
		return exprToString(t.X) + "." + t.Sel.Name
	// case *ast.StructType: // 匿名结构体
	// case *ast.InterfaceType: // 接口类型
	// case *ast.FuncType: // 函数类型
	// case *ast.ChanType: // Channel 类型
	default:
		// 对于更复杂的类型，可以返回一个占位符或进行更详细的处理
		return fmt.Sprintf("%T", expr)
	}
}

// parseArrayLength 解析数组长度（处理字面量和常量）
func parseArrayLength(lenExpr ast.Expr) (int64, error) {
	switch expr := lenExpr.(type) {
	case *ast.BasicLit:
		// 字面量类型（如[5]int中的5）
		if expr.Kind == token.INT {
			return strconv.ParseInt(expr.Value, 10, 64)
		}
		return -1, fmt.Errorf("不支持的数组长度字面量类型：%s", expr.Kind)
	case *ast.Ident:
	default:
		return -1, fmt.Errorf("不支持的数组长度表达式类型：%T", lenExpr)
	}
	return -1, fmt.Errorf("未找到数组长度")
}

var (
	hf_plat  *string
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
	hf_plat = handlerCmd.Flags().StringP("platform", "p", "", "选择平台")
}
