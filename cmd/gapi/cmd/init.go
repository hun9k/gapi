/*
Copyright © 2025 9k <hun9k.github.io>
The MIT License (MIT)
*/
package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"os/exec"
	"path"

	"github.com/hun9k/gapi/internal/tmpls"
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
		// module-path is required
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return errors.New("module-path is required")
		}
		// check module-path
		// if !utils.IsValidModulePath(args[0]) {
		// 	return errors.New("module-path is invalid")
		// }
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// // work dir
		// wd, err := os.Getwd()
		// if err != nil {
		// 	slog.Error("get work dir was failed", "error", err)
		// 	return
		// }

		// mod info
		mod := modInfo{
			modPath: args[0],
			modBase: path.Base(args[0]),
		}

		// generate basic structure
		if err := genBasicStructure(mod); err != nil {
			slog.Error("generate basic structure was failed", "error", err)
			return
		}

		// change dir to mod dir
		if err := changeDir(mod); err != nil {
			slog.Error("create and change dir was failed", "error", err)
			return
		}

		// mod init
		if err := modInit(mod); err != nil {
			slog.Error("mod init was failed", "error", err)
			return
		}

		// generate basic codes
		if err := genBasicCodes(mod); err != nil {
			slog.Error("generate basic codes was failed", "error", err)
			return
		}

		slog.Info("new module initialized", "module-path", mod.modPath)
		slog.Info("please execute:")
		slog.Info(fmt.Sprintf("cd %s", mod.modBase))
		slog.Info("go mod tidy")
		slog.Info("go run .")
	},
}

type modInfo struct {
	modPath, modBase string
}

// create module dir and change work dir to it
func changeDir(mod modInfo) error {
	// change work dir to module path
	if err := os.Chdir(mod.modBase); err != nil {
		return err
	}

	return nil
}

// mod init
func modInit(mod modInfo) error {
	// execute `go mod init`
	cmdModInit := exec.Command("go", "mod", "init", mod.modPath)
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

type modFile struct {
	isDir    bool
	filename string
	content  string
}

// generate basic structure
func genBasicStructure(mod modInfo) error {
	// structured files
	files := [...]modFile{
		{true, mod.modBase, ""},
		{false, "README.md", fmt.Sprintf("# %s\n", mod.modPath)},
		{true, "apis", ""},
		{false, "apis/README.md", "# your APIs in here\n"},
		{true, "internal", ""},
		{true, "internal/schemas", ""},
		{false, "internal/schemas/README.md", "# your schemas in here\n"},
	}
	const dirMode, fileMode = 0754, 0640
	for _, f := range files {
		if f.isDir {
			if err := os.Mkdir(f.filename, dirMode); err != nil {
				return err
			}
		} else {
			if err := os.WriteFile(f.filename, []byte(f.content), fileMode); err != nil {
				return err
			}
		}
	}

	return nil
}

type codeTmpl struct {
	name, text, filename string
}

// generate basic codes
func genBasicCodes(mod modInfo) error {
	codeTmpls := [...]codeTmpl{
		{"main", tmpls.Main, "main.go"},
		{"api_groups", tmpls.Apis_group, "apis/groups.go"},
		{"main", tmpls.Main, "main.go"},
		{"main", tmpls.Main, "main.go"},
	}

	for _, ct := range codeTmpls {
		tmpl := template.Must(template.New(ct.name).Parse(ct.text))
		file, err := os.Create(ct.filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := tmpl.Execute(file, mod); err != nil {
			return err
		}
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
