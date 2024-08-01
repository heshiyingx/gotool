func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) {
    {{if .withCache}}
        {{.cacheKey}}

		var {{.upperStartCamelPrimaryKey}} {{.pkNameType}}
        err := m.db.QueryToGetPKCtx(ctx, {{.cacheKeyVariable}}, &{{.upperStartCamelPrimaryKey}}, func(ctx context.Context, p *{{.pkNameType}}, db *gorm.DB) error {
            return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Select("{{.pkNameWrap}}").Where("{{.originalField}}", {{.lowerStartCamelField}}).Take(p).Error
        })
        if err != nil {
            return nil, err
        }
        {{.pkKey}}
        var resp {{.upperStartCamelObject}}
		err = m.db.QueryOneByPKCtx(ctx, &resp, {{.pkCacheKeyName}}, func(ctx context.Context, r any, db *gorm.DB) error {
            return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}} = ?", {{.upperStartCamelPrimaryKey}}).Take(r).Error
        })
        return &resp, err

    {{else}}
        var resp Users
        err := m.db.QueryCtx(ctx, &resp, func(ctx context.Context, r any, db *gorm.DB) error {
        return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalField}}", {{.lowerStartCamelField}}).Take(r).Error
        })
        return &resp, err
    {{end}}
}

func (m *default{{.upperStartCamelObject}}Model) UpdateOneBy{{.upperField}}(ctx context.Context, {{.in}},updateObj *{{.upperStartCamelObject}},delCacheKeys []string,fields ...string) (int64, error) {
    if updateObj==nil{
		return 0,nil
    }
    {{if .withCache}}
        data,err := m.FindOneBy{{.upperField}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
            return 0, err
		}
		{{.keys}}

        delKeys := []string{ {{.keyNames}} }
        if len(delCacheKeys) > 0 { delKeys = append(delKeys, delCacheKeys...) }

        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}}", {{.lowerStartCamelField}})
            if len(fields) > 0 {
                upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
            }else {
                upTx = upTx.Save(updateObj)
            }
            return upTx.RowsAffected,upTx.Error
        },delKeys...)

    {{else}}

    return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
        upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalField}}", {{.lowerStartCamelField}})
        if len(fields) > 0 {
            upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
        } else {
            upTx = upTx.Save(updateObj)
        }
        return upTx.RowsAffected, upTx.Error
    })
    {{end}}
}
func (m *default{{.upperStartCamelObject}}Model) DeleteOneBy{{.upperField}}(ctx context.Context, {{.in}},delCacheKeys ...string) (int64, error) {

    {{if .withCache}}
        data,err := m.FindOneBy{{.upperField}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
        return 0, err
        }
        {{.keys}}

        delKeys := []string{ {{.keyNames}} }
        if len(delCacheKeys) > 0 { delKeys = append(delKeys, delCacheKeys...) }

        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            delTx := db.WithContext(ctx).Where("{{.pkNameWrap}}", {{.lowerStartCamelField}}).Delete(&{{.upperStartCamelObject}}{})
            return delTx.RowsAffected,delTx.Error
        },delKeys...)

    {{else}}
    return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
        delTx := db.WithContext(ctx).Where("{{.originalField}}", {{.lowerStartCamelField}}).Delete(&Users{})
        return delTx.RowsAffected, delTx.Error
    })
    {{end}}
    }