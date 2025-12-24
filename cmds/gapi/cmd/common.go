package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hun9k/gapi/cmds/gapi/internal/tmpls"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

const ()

const (
	HANDLER_DIR        = "handlers"
	TASK_DIR           = "tasks"
	MODEL_DIR          = "models"
	GO_EXT             = ".go"
	LOCK_FILE          = ".gapi.lock"
	CONFIG_FILE        = "configs.yaml"
	README_FILE        = "README.md"
	DOCKERFILE_FILE    = "Dockerfile"
	DOCKERCOMPOSE_FILE = "docker-compose.yaml"
	GITIGNORE_FILE     = ".gitignore"
	MAIN_FILE          = "main.go"
	INIT_FILE          = "init.go"
	MESSAGE_FILE       = "messsage_gen.go"
	HANDLER_FILE       = "handlers_gen.go"
	ROUTER_FILE        = "routers_gen.go"
	SETUP_FILE         = "setup.go"
)

type fileTmpl struct {
	text, filename string
	data           any
	isKeep         bool // 保持不变
	isAlways       bool // 每次更新
}

const (
	DIR_MODE  = 0754
	FILE_MODE = 0640
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func MkDir(dir string) error {
	if err := os.MkdirAll(dir, DIR_MODE); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// generate structure
func MkDirs(dirs []string) []error {
	errs := make([]error, len(dirs))
	for i, dir := range dirs {
		errs[i] = MkDir(dir)
	}
	return errs
}

func GenFile(ft fileTmpl) error {
	// 文件存在 同时 保留或不强制覆盖
	if _, err := os.Stat(ft.filename); err == nil && (ft.isKeep || !*flagForce) && !ft.isAlways {
		return fmt.Errorf("文件未更新")
	}

	t, err := template.New("").Parse(ft.text)
	if err != nil {
		return err
	}

	src := new(bytes.Buffer)
	if err := t.Execute(src, ft.data); err != nil {
		return err
	}

	// format go content
	content := src.Bytes()
	if filepath.Ext(ft.filename) == GO_EXT {
		c, err := format.Source(src.Bytes())
		if err != nil {
			return err
		}
		content = c
	}

	// write
	if err := os.WriteFile(ft.filename, content, FILE_MODE); err != nil {
		return err
	}

	return nil
}

// generate codes
func GenFiles(tmpls []fileTmpl) []error {
	errs := make([]error, len(tmpls))
	for i, ct := range tmpls {
		errs[i] = GenFile(ct)
	}
	return errs
}

func ModFile(filename string) (*modfile.File, error) {
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

// 生成资源模型文件
func GenModel(res string) error {
	tmpl := fileTmpl{
		text:     tmpls.ResModel,
		filename: filepath.Join(MODEL_DIR, res+GO_EXT),
		data: tmplData{
			"package":  MODEL_DIR,
			"model":    strcase.ToCamel(res),
			"resource": res,
		},
	}
	if err := GenFile(tmpl); err != nil {
		slog.Error("生成模型失败", "filename", tmpl.filename, "error", err)
		return err
	}
	slog.Info("生成模型完成，应基于业务编辑模型", "filename", tmpl.filename)

	return nil
}

// 更新models/init.go
func SetModelInit(modelPath string) error {
	// 获取模型列表
	modelList, err := getModels(modelPath)
	if err != nil {
		slog.Error("模型列表失败", "error", err)
		return err
	}
	slog.Info("模型列表成功", "modelList", strings.Join(modelList, ", "))

	// 更新 models/init.go
	initTmpl := fileTmpl{
		text:     tmpls.ModelsInit,
		filename: filepath.Join(modelPath, "init.go"),
		data: tmplData{
			"modelList": strings.Join(modelList, ", "),
		},
		isAlways: true,
	}
	if err := GenFile(initTmpl); err != nil {
		slog.Error("初始模型失败", "filename", initTmpl.filename, "error", err)
		return err
	}
	slog.Info("初始模型成功", "filename", initTmpl.filename)

	return nil
}

func GenHandler(appPath string, mod *modfile.File, resource string) error {
	// 解析plat/resName, 例如admin/post
	plat, resName := parseResource(resource)

	// 平台目录
	platPath := filepath.Join(appPath, HANDLER_DIR, plat)

	// 创建资源目录
	resPath := filepath.Join(platPath, resName)
	if err := MkDir(resPath); err != nil {
		slog.Error("创建目录失败", "filename", resPath, "error", err)
		return err
	}
	slog.Info("创建目录成功", "filename", resPath)

	// 生成资源相关文件
	fields, err := ParseModel(resName)
	if err != nil {
		slog.Error("解析模型失败", "error", err)
		return err
	}
	slog.Info("解析模型成功", "filename", filepath.Join(appPath, MODEL_DIR, resName+".go"))

	// 生成messages，handlers，routers，setup文件
	modelName := "models." + strcase.ToCamel(resName)
	codeTmpls := []fileTmpl{
		{
			text:     tmpls.ResMessages,
			filename: filepath.Join(resPath, "messages_gen.go"),
			data: tmplData{
				"resource":   resName,
				"fields":     fields,
				"modelName":  modelName,
				"modPath":    mod.Module.Mod.Path,
				"iTime":      iTime,
				"iSql":       iSql,
				"iDatatypes": iDatatypes,
			},
		},
		{
			text:     tmpls.ResHandlers,
			filename: filepath.Join(resPath, "handlers_gen.go"),
			data: tmplData{
				"resource":  resName,
				"modelName": modelName,
				"modPath":   mod.Module.Mod.Path,
			},
		},
		{
			text:     tmpls.ResRouters,
			filename: filepath.Join(resPath, "routers_gen.go"),
			data: tmplData{
				"resource": resName,
			},
		},
		{
			text:     tmpls.ResSetup,
			filename: filepath.Join(resPath, SETUP_FILE),
			data: tmplData{
				"resource":     resName,
				"routerPrefix": resName,
			},
			isKeep: true,
		},
	}
	for i, err := range GenFiles(codeTmpls) {
		if err != nil {
			slog.Error("生成代码失败", "filename", codeTmpls[i].filename, "error", err)
			continue
		}
		slog.Info("生成代码成功", "filename", codeTmpls[i].filename)
	}

	return nil
}

func SetPlat(appPath string, mod *modfile.File, resource string) error {
	// 解析plat/resName, 例如admin/post
	plat, res := parseResource(resource)
	user := ""
	if res == "user" {
		user = "user"
	}

	// 平台代码
	pSetup, pRouter := SETUP_FILE, ROUTER_FILE
	if plat != "" {
		pSetup = plat + "_" + pSetup
	}
	if plat != "" {
		pRouter = plat + "_" + pRouter
	}
	codeTmpls := []fileTmpl{
		{
			text:     tmpls.PlatSetup,
			filename: filepath.Join(appPath, HANDLER_DIR, pSetup),
			data: tmplData{
				"platform": plat,
				"user":     user,
			},
			isKeep: true,
		},
		{
			text:     tmpls.PlatRouters,
			filename: filepath.Join(appPath, HANDLER_DIR, pRouter),
			data: tmplData{
				"platform":  plat,
				"modPath":   mod.Module.Mod.Path,
				"resources": getPlatRes(filepath.Join(HANDLER_DIR, plat)),
			},
		},
	}
	for i, err := range GenFiles(codeTmpls) {
		if err != nil {
			slog.Error("生成文件失败", "filename", codeTmpls[i].filename, "error", err)
			return err
		}
		slog.Info("生成文件成功", "filename", codeTmpls[i].filename)
	}

	return nil
}

// 生成用户模型和Handler
func GenUserModelHandler(appPath string, mod *modfile.File, plat string) error {
	res, resource := "user", plat+"/user"

	// 生成用户模型文件
	modelTmpl := fileTmpl{
		text:     tmpls.UserModel,
		filename: filepath.Join(appPath, MODEL_DIR, res+".go"),
		data: tmplData{
			"package": MODEL_DIR,
		},
	}
	if err := GenFile(modelTmpl); err != nil {
		slog.Error("生成文件失败", "filename", modelTmpl.filename, "error", err)
		return err
	}
	slog.Info("生成文件成功", "filename", modelTmpl.filename)

	// 更新模型init
	if err := SetModelInit(filepath.Join(appPath, MODEL_DIR)); err != nil {
		return err
	}

	// 生成用户通用Handler文件
	if err := GenHandler(appPath, mod, resource); err != nil {
		return err
	}

	// 设置平台相关代码
	if err := SetPlat(appPath, mod, resource); err != nil {
		return err
	}

	// 生成用户特定文件
	codeTmpls := []fileTmpl{
		{
			text:     tmpls.UserMessage,
			filename: filepath.Join(appPath, HANDLER_DIR, plat, res, "messages.go"),
			data: tmplData{
				"resource": res,
			},
		},
		{
			text:     tmpls.UserHandler,
			filename: filepath.Join(appPath, HANDLER_DIR, plat, res, "handlers.go"),
			data: tmplData{
				"modPath":   mod.Module.Mod.Path,
				"resource":  res,
				"modelName": "models." + strcase.ToCamel(res),
			},
		},
		{
			text:     tmpls.UserRouter,
			filename: filepath.Join(appPath, HANDLER_DIR, plat, res, "routers.go"),
			data: tmplData{
				"resource": res,
			},
		},
	}
	for i, err := range GenFiles(codeTmpls) {
		if err != nil {
			slog.Error("生成文件失败", "filename", codeTmpls[i].filename, "error", err)
			return err
		}
		slog.Info("生成文件成功", "filename", codeTmpls[i].filename)
	}

	return nil
}

func parseResource(res string) (plat string, resource string) {
	strs := strings.Split(res, "/")
	if len(strs) == 1 {
		return "", strs[0]
	} else {
		return strs[0], strs[1]
	}
}
