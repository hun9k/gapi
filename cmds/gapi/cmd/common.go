package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"go.yaml.in/yaml/v4"
	"golang.org/x/mod/modfile"
)

const (
	GO_EXT          = ".go"
	OPENAPI_VERSION = "3.0.4"
)

type codeTmpl struct {
	text, filename string
	data           any
}

const (
	DIR_MODE  = 0754
	FILE_MODE = 0640
)

func genOpenApi() ([]byte, error) {
	// schema
	schema := openapi3.NewSchema()
	schema.Type = &openapi3.Types{"object"}
	schema.Properties = map[string]*openapi3.SchemaRef{
		"StringField": {
			Value: &openapi3.Schema{
				Type:     &openapi3.Types{"string"},
				Title:    "字符串类型",
				Nullable: false,
			},
		},
		"NullStringField": {
			Value: &openapi3.Schema{
				Type:     &openapi3.Types{"string"},
				Title:    "可空字符串类型",
				Nullable: true,
			},
		},
	}

	// schemas
	schemas := openapi3.Schemas{
		"contents": schema.NewRef(),
	}

	// components
	components := openapi3.Components{
		Schemas: schemas,
	}

	// openapi
	doc := openapi3.T{
		OpenAPI: OPENAPI_VERSION,
	}
	doc.Info = &openapi3.Info{
		Title:   "gapi",
		Version: "0.0.1",
	}
	doc.Paths = &openapi3.Paths{}
	doc.Components = &components

	// c, err := doc.MarshalYAML()
	// if err != nil {
	// 	return nil, err
	// }

	out, err := yaml.Marshal(doc)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(out))

	return out, nil
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
