package cmd

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	GO_EXT = ".go"
)

type codeTmpl struct {
	text, filename string
	data           any
}

type modInfo struct {
	Path, Base, Dir string
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

type dir struct {
	filename string
}

// generate structure
func genDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := os.Mkdir(dir, DIR_MODE); err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}

// generate codes
func genCodes(tmpls []codeTmpl) error {
	for _, ct := range tmpls {
		tmpl, err := template.New("").Parse(ct.text)
		if err != nil {
			return err
		}

		src := new(bytes.Buffer)
		if err := tmpl.Execute(src, ct.data); err != nil {
			return err
		}

		// format go code
		code := src.Bytes()
		if filepath.Ext(ct.filename) == GO_EXT {
			c, _ := format.Source(src.Bytes())
			code = c
		}

		// write
		if err := os.WriteFile(ct.filename, code, FILE_MODE); err != nil {
			return err
		}
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
