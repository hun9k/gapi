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

	"github.com/hun9k/gapi"
	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init module-path",
	Short: "initialize an module's structure",
	Long: `For example:
gapi init github.com/hun9k/gapi-module

A module will be generated in gapi-module directory, and the module-path will be github.com/hun9k/gapi-module.`,
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
			{false, filepath.Join(mod.Base, ROUTER_BASE, README_BASE), "# your routers in here\n"},
			{true, filepath.Join(mod.Base, INTERNAL_BASE), ""},
			{true, filepath.Join(mod.Base, INTERNAL_BASE, SCHEMAS_BASE), ""},
			{false, filepath.Join(mod.Base, INTERNAL_BASE, SCHEMAS_BASE, README_BASE), "# your schemas in here\n"},
		}
		if err := genStructure(basicFiles); err != nil {
			slog.Error("generate basic structure failed", "error", err)
			return
		}

		// change dir to mod dir
		if err := changeDir(mod); err != nil {
			slog.Error("create and change dir failed", "error", err)
			return
		}

		// mod init
		if err := modInit(mod); err != nil {
			slog.Error("mod init failed", "error", err)
			return
		}

		// generate basic codes
		codeTmpls := []codeTmpl{
			{tmpls.Configs, CONFIGS_BASE, gapi.NewDefaultConf()},
			{tmpls.Apis_group, filepath.Join(ROUTER_BASE, GROUP_ROUTER_BASE), mod},
			{tmpls.Main, MAIN_BASE, mod},
		}
		if err := genCodes(codeTmpls); err != nil {
			slog.Error("generate basic codes failed", "error", err)
			return
		}

		// mod tidy
		if err := modTidy(mod); err != nil {
			slog.Error("mod tidy failed", "error", err)
			return
		}

		// success and tips
		slog.Info("new module initialized", "module-path", mod.Path)
		slog.Info("you should modify configs", "file", "configs.yaml")
		slog.Info("and exec", "cmd", fmt.Sprintf("cd %s && go run .", mod.Base))
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
	cmdModReplace := exec.Command("go", "mod", "edit", `-replace=github.com/hun9k/gapi=D:/apps/gapi`)
	if err := cmdModReplace.Run(); err != nil {
		return err
	}
	return nil
}

// mod tidy
func modTidy(mod modInfo) error {
	// execute `go mod tidy`
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
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
