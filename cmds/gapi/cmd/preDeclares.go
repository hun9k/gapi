package cmd

import (
	"html/template"
	"os"
)

const (
	ROUTER_BASE       = "routers"
	INTERNAL_BASE     = "internal"
	SCHEMAS_BASE      = "schemas"
	README_BASE       = "README.md"
	CONFIGS_BASE      = "configs.yaml"
	MAIN_BASE         = "main.go"
	GROUP_ROUTER_BASE = "groups.go"
)

type codeTmpl struct {
	text, filename string
	data           any
}

type modInfo struct {
	Path, Base string
}

type resourceInfo struct {
	Mod      modInfo
	Resource string
	Version  string
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
			if err := os.Mkdir(f.filename, DIR_MODE); err != nil {
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
	}

	return nil
}
