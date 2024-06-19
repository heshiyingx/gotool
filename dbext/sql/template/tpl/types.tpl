type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}
			db   *gormdb.CacheGormDB[{{.upperStartCamelObject}}, {{.pkType}}]
			redisClient *redis.Client
		{{else}}
			{{/*conn sqlx.SqlConn*/}}
		{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
