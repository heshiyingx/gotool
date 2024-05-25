type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}
			db   *gormdb.CacheGormDB[{{.upperStartCamelObject}}, {{.pkType}}]
		{{else}}
			db   *gormdb.GormDB[{{.upperStartCamelObject}}, {{.pkType}}]
		{{end}}
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
