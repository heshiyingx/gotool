import (
	"context"
	{{if .table.ContainsNullField}}"database/sql"{{end}}
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

    "github.com/heshiyingx/gotool/dbext/gormdb"
    "gorm.io/gorm"

	{{.third}}
)
