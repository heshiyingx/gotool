import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

    "github.com/heshiyingx/gotool/dbext/gormdb"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"

	{{.third}}
)
