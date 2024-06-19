func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.db.QueryOneByPKCtx(ctx, &resp,  {{.cacheKeyVariable}}, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.Model(&upperStartCamelObject{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}}).Take(r).Error
	})
	return &resp,err
	}{{end}}
}
