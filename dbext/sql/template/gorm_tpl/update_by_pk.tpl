func (m *default{{.upperStartCamelObject}}Model) UpdateBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},updateObj *{{.upperStartCamelObject}},delCacheKeys []string,fields ...string) (int64, error) {
	if updateObj==nil{
		return 0,nil
	}
	{{if .withCache}}
	data,err := m.FindBy{{.titlePrimaryKey}}(ctx,{{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return 0, err
	}
	{{.keys}}

	delKeys := []string{ {{.keyNames}} }
	if len(delCacheKeys) > 0 { delKeys = append(delKeys, delCacheKeys...) }

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}})
		if len(fields) > 0 {
			upTx = upTx.Select(fields).Updates(updateObj)
		}else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected,upTx.Error
	},delKeys...)
	{{else}}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}})
	if len(fields) > 0 {
		upTx = upTx.Select(fields).Updates(updateObj)
	} else {
		upTx = upTx.Save(updateObj)
	}
		return upTx.RowsAffected, upTx.Error
	})

	{{end}}
}
