func (m *default{{.upperStartCamelObject}}Model) DeleteBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (int64, error) {
    {{if .withCache}}
        {{if .containsIndexCache}}
            data, err := m.FindBy{{.titlePrimaryKey}}(ctx, {{.lowerStartCamelPrimaryKey}})
            if err != nil {
                return 0,err
            }
        {{end}}
        {{.keys}}

        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
		}, {{.keyValues}})
    {{else}}
        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
        })
    {{end}}
}
