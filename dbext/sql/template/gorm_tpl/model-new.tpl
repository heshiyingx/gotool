func newDefault{{.upperStartCamelObject}}Model({{if .withCache}}config gormdb.Config{{end}}) *default{{.upperStartCamelObject}}Model {

    {{if .withCache}}
        cacheGormDB := gormdb.MustNewCacheGormDB[{{.upperStartCamelObject}}, {{.pkType}}](config)
        return &default{{.upperStartCamelObject}}Model{
        db: cacheGormDB,
        table:      {{.table}},
        }
    {{else}}
    {{end}}
}

