/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

// genschemaCmd represents the genschema command
var genschemaCmd = &cobra.Command{
	Use:   "genschema resource",
	Short: "生成资源设计",
	Long: `生成用于GORM使用的资源模型。例如：
	gapi genschema contents
会生成contents的GORM模型：
	internal/schemas/contents.go
等代码`,
	Args: func(cmd *cobra.Command, args []string) error {
		// 必须指定至少一个resource
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, resources []string) {
		slog.Info("生成资源资源结构代码开始", "resources", resources, "--hooks", *genschemaHooks)

		// get module info
		modFile, err := modFileByFile()
		if err != nil {
			slog.Error("获取module信息失败", "error", err)
			return
		}
		mod := modInfo{
			Path: modFile.Module.Mod.Path,
		}

		hs := parseHooks(*genschemaHooks)
		for _, resource := range resources {
			rInfo := resourceInfo{
				Mod:      mod,
				Resource: resource,
				Schema: schemaInfo{
					Name:  resourceSchemaName(resource),
					Model: *genSchemaModel,
					Hooks: hs,
				},
			}

			// codeTmpls
			codeTmpls := []codeTmpl{
				// {tmpls.Resource_schemas, filepath.Join(INTERNAL_BASE, SCHEMAS_BASE, resource+GO_EXT), rInfo},
			}

			// generate codes
			slog.Info("生成资源结构代码", "resource", rInfo.Resource)
			if err := genCodes(codeTmpls); err != nil {
				slog.Error("生成资源结构代码失败", "error", err)
				return
			}
		}
		slog.Info("生成资源资源结构代码完成", "resources", resources)
		slog.Info("接下来应自定义资源结构，然后生成API代码")

	},
}

func parseHooks(hook string) map[string]bool {
	hooks := map[string]bool{
		"save": false, "create": false, "update": false, "delete": false, "find": false,
	}
	for _, h := range strings.Split(hook, ",") {
		h = strings.TrimSpace(h)
		if _, exists := hooks[h]; !exists && h != "all" {
			continue
		}

		if h == "all" {
			hooks["save"], hooks["create"], hooks["update"], hooks["delete"], hooks["find"] = true, true, true, true, true
			break
		} else {
			hooks[h] = true
		}
	}
	return hooks
}

// flags
var (
	genschemaHooks *string
	genSchemaModel *bool
)

func init() {
	rootCmd.AddCommand(genschemaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genschemaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genschemaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genschemaHooks = genschemaCmd.Flags().StringP("hooks", "k", "", "默认不包含钩子，create,update,delete,find,all表示某类或全部")
	genSchemaModel = genschemaCmd.Flags().BoolP("model", "m", true, "默认true表示嵌入Model，false表示不嵌入")
}
