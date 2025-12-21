package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/mod/modfile"
)

const (
	GO_EXT    = ".go"
	MODEL_DIR = "models"
	LOCK_FILE = ".gapi.lock"
)

type codeTmpl struct {
	text, filename string
	data           any
}

const (
	DIR_MODE  = 0754
	FILE_MODE = 0640
)

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// generate structure
func genDirs(dirs []string) []error {
	errs := make([]error, len(dirs))
	for i, dir := range dirs {
		if err := os.MkdirAll(dir, DIR_MODE); err != nil && !os.IsExist(err) {
			errs[i] = err
			continue
		}
	}
	return errs
}

// generate codes
func genCodes(tmpls []codeTmpl, force bool) []error {
	errs := make([]error, len(tmpls))
	for i, ct := range tmpls {
		if !force {
			if _, err := os.Stat(ct.filename); err == nil {
				errs[i] = fmt.Errorf("文件已存在，-f强制生成")
				continue
			}
		}

		tmpl, err := template.New("").Parse(ct.text)
		if err != nil {
			errs[i] = err
			continue
		}

		src := new(bytes.Buffer)
		if err := tmpl.Execute(src, ct.data); err != nil {
			errs[i] = err
			continue
		}

		// format go code
		code := src.Bytes()
		if filepath.Ext(ct.filename) == GO_EXT {
			c, err := format.Source(src.Bytes())
			if err != nil {
				errs[i] = err
				continue
			}
			code = c
		}

		// write
		if err := os.WriteFile(ct.filename, code, FILE_MODE); err != nil {
			errs[i] = err
			continue
		}
	}
	return errs
}

func modFile(filename string) (*modfile.File, error) {
	// 读取 go.mod 文件内容
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析 go.mod 内容
	modFile, err := modfile.Parse(filename, modBytes, nil)
	if err != nil {
		return nil, err
	}

	return modFile, nil
}

type tmplData map[string]any
