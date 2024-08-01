func (m *default{{.upperStartCamelObject}}Model) DeleteBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},delCacheKeys ...string) (int64, error) {
    {{if .withCache}}
        {{if .containsIndexCache}}
            data, err := m.FindBy{{.titlePrimaryKey}}(ctx, {{.lowerStartCamelPrimaryKey}})
            if err != nil {
                return 0,err
            }
        {{end}}
        {{.keys}}
        delKeys := []string{ {{.cacheKeyStr}} }
		if len(delCacheKeys) > 0 { delKeys = append(delKeys, delCacheKeys...) }


        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.WithContext(ctx).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
		}, delKeys...)
    {{else}}
        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.WithContext(ctx).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
        })
    {{end}}
}
