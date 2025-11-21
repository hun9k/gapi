package cmd

import (
	"errors"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ROUTER_BASE       = "routers"
	INTERNAL_BASE     = "internal"
	SCHEMAS_BASE      = "schemas"
	README_BASE       = "README.md"
	CONFIGS_BASE      = "configs.yaml"
	GO_EXT            = ".go"
	MAIN_BASE         = "main.go"
	GROUP_ROUTER_BASE = "groups.go"
	BIZS_BASE         = "bizs.go"
	HANDLERS_BASE     = "handlers.go"
)

type codeTmpl struct {
	text, filename string
	data           any
}

type modInfo struct {
	Path, Base string
}

type resourceInfo struct {
	Mod               modInfo
	Version           string
	Resource          string
	Schema            schemaInfo
	ResourceBody      string
	ResourcePatchBody string
}

type schemaInfo struct {
	Name  string
	Model bool
	Hooks map[string]bool
}

func resourceSchemaName(resource string) string {
	return cases.Title(language.English).String(resource)
}

const (
	DIR_MODE  = 0754
	FILE_MODE = 0640
)

type structureFile struct {
	isDir    bool
	filename string
	content  string
}

// generate structure
func genStructure(files []structureFile) error {

	for _, f := range files {
		if f.isDir {
			if err := os.Mkdir(f.filename, DIR_MODE); err != nil && !os.IsExist(err) {
				return err
			}
		} else {
			if err := os.WriteFile(f.filename, []byte(f.content), FILE_MODE); err != nil {
				return err
			}
		}
	}

	return nil
}

// generate codes
func genCodes(tmpls []codeTmpl) error {
	for _, ct := range tmpls {
		tmpl := template.Must(template.New(ct.filename).Parse(ct.text))
		file, err := os.Create(ct.filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := tmpl.Execute(file, ct.data); err != nil {
			return err
		}

		// format code
		fmtCode(ct.filename)
	}

	return nil
}

func fmtCode(filename string) error {
	if filepath.Ext(filename) != GO_EXT {
		return errors.New("non .go file")
	}

	cmd := exec.Command("go", "fmt", filename)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func modFileByFile() (*modfile.File, error) {
	goModFilename := "go.mod"
	// 读取 go.mod 文件内容
	modBytes, err := os.ReadFile(goModFilename)
	if err != nil {
		return nil, err
	}

	// 解析 go.mod 内容
	modFile, err := modfile.Parse(goModFilename, modBytes, nil)
	if err != nil {
		return nil, err
	}

	return modFile, nil
}
