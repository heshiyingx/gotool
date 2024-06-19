func newDefault{{.upperStartCamelObject}}Model(config gormdb.Config) *default{{.upperStartCamelObject}}Model {

    {{if .withCache}}
        cacheGormDB := gormdb.MustNewCacheGormDB[{{.upperStartCamelObject}}, {{.pkType}}](config)
        return &default{{.upperStartCamelObject}}Model{
            db: cacheGormDB,
        }
    {{else}}
    gormDB := gormdb.MustNewGormDB[{{.upperStartCamelObject}}, {{.pkType}}](config)
	return &default{{.upperStartCamelObject}}Model{
		db: gormDB,
    }
    {{end}}
}

