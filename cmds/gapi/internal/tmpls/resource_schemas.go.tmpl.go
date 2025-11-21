package tmpls

var Resource_schemas = `package schemas
{{$save := index .Schema.Hooks "save"}}
{{$create := index .Schema.Hooks "create"}}
{{$update := index .Schema.Hooks "update"}}
{{$delete := index .Schema.Hooks "delete"}}
{{$find := index .Schema.Hooks "find"}}
{{if or .Schema.Model $save $create $update $delete $find}}
import (
{{if .Schema.Model}}	"github.com/hun9k/gapi"{{end}}
{{if or $save $create $update $delete $find}}	"gorm.io/gorm"{{end}}
)
{{end}}
type {{.Schema.Name}} struct {
	// 自定义字段

{{if .Schema.Model}}	// 嵌入基础字段，ID，CreatedAt，UpdatedAt, DeletedAt
	gapi.Model
{{end}}
}
{{if or $save $create $update $delete $find}}
// hooks
{{end}}
{{if $save}}
// save order: beforeSave, create or update, afterSave
func (m *{{.Schema.Name}}) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (m *{{.Schema.Name}}) AfterSave(tx *gorm.DB) error {
	return nil
}
{{end}}
{{if $create}}
// create
func (m *{{.Schema.Name}}) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (m *{{.Schema.Name}}) AfterCreate(tx *gorm.DB) error {
	return nil
}
{{end}}
{{if $update}}
// update
func (m *{{.Schema.Name}}) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (m *{{.Schema.Name}}) AfterUpdate(tx *gorm.DB) error {
	return nil
}
{{end}}
{{if index .Schema.Hooks "delete"}}
// delete
func (m *{{.Schema.Name}}) BeforeDelete(tx *gorm.DB) error {
	return nil
}

func (m *{{.Schema.Name}}) AfterDelete(tx *gorm.DB) error {
	return nil
}
{{end}}
{{if index .Schema.Hooks "find"}}
// find
func (m *{{.Schema.Name}}) AfterFind(tx *gorm.DB) error {
	return nil
}
{{end}}
`
