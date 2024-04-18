package {{.pkg}}

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

var(
  {{if .withCache}}{{.cacheKeys}}{{end}}
)
