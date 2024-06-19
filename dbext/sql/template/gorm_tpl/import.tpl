import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

    "github.com/heshiyingx/gotool/dbext/gormdb"
    "gorm.io/gorm"

	{{.third}}
)
