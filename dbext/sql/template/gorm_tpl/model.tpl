package {{.pkg}}
{{if .withCache}}
import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)
{{else}}


{{end}}
var _ {{.upperStartCamelObject}}DBInterface = (*custom{{.upperStartCamelObject}}DB)(nil)

type (
	// {{.upperStartCamelObject}}DBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}DBInterface interface {
		{{.lowerStartCamelObject}}Model
	}

	custom{{.upperStartCamelObject}}DB struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}DB returns a model for the database table.
func New{{.upperStartCamelObject}}DB(config gormdb.Config) {{.upperStartCamelObject}}DBInterface {
	return &custom{{.upperStartCamelObject}}DB{
		default{{.upperStartCamelObject}}Model: newDefault{{.upperStartCamelObject}}Model(config),
	}
}

{{if not .withCache}}

{{end}}

