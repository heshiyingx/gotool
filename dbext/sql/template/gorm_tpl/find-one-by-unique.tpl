func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) {
    {{if .withCache}}
        {{.cacheKey}}

		var {{.upperStartCamelPrimaryKey}} {{.pkNameType}}
        err := m.db.QueryToGetPKCtx(ctx, {{.cacheKeyVariable}}, &{{.upperStartCamelPrimaryKey}}, func(ctx context.Context, p *{{.pkNameType}}, db *gorm.DB) error {
            return db.Model(&{{.upperStartCamelObject}}{}).Select("{{.pkNameWrap}}").Where("{{.originalField}}", {{.lowerStartCamelField}}).Take(p).Error
        })
        if err != nil {
            return nil, err
        }
        {{.pkKey}}
        var resp {{.upperStartCamelObject}}
		err = m.db.QueryOneByPKCtx(ctx, &resp, {{.pkCacheKeyName}}, func(ctx context.Context, r any, db *gorm.DB) error {
            return db.Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}} = ?", {{.upperStartCamelPrimaryKey}}).Take(r).Error
        })
        return &resp, err

    {{else}}
        {{/*var resp {{.upperStartCamelObject}}*/}}
        {{/*query := fmt.Sprintf("select %s from %s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)*/}}
        {{/*err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelField}})*/}}
        {{/*switch err {*/}}
        {{/*case nil:*/}}
        {{/*    return &resp, nil*/}}
        {{/*case sqlx.ErrNotFound:*/}}
        {{/*    return nil, ErrNotFound*/}}
        {{/*default:*/}}
        {{/*    return nil, err*/}}
        {{/*}*/}}
    {{end}}
}

func (m *default{{.upperStartCamelObject}}Model) UpdateOneBy{{.upperField}}(ctx context.Context, {{.in}},updateObj *{{.upperStartCamelObject}},fields ...string) (int64, error) {
    if updateObj==nil{
		return 0,nil
    }
    {{if .withCache}}
        data,err := m.FindOneBy{{.upperField}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
            return 0, err
		}
		{{.keys}}
        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            upTx := db.Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}}", {{.lowerStartCamelField}})
            if len(fields) > 0 {
                upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
            }else {
                upTx = upTx.Save(updateObj)
            }
            return upTx.RowsAffected,upTx.Error
        },{{.keyNames}})

    {{else}}
        {{/*var resp {{.upperStartCamelObject}}*/}}
        {{/*query := fmt.Sprintf("select %s from %s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)*/}}
        {{/*err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelField}})*/}}
        {{/*switch err {*/}}
        {{/*case nil:*/}}
        {{/*    return &resp, nil*/}}
        {{/*case sqlx.ErrNotFound:*/}}
        {{/*    return nil, ErrNotFound*/}}
        {{/*default:*/}}
        {{/*    return nil, err*/}}
        {{/*}*/}}
    {{end}}
}
func (m *default{{.upperStartCamelObject}}Model) DeleteOneBy{{.upperField}}(ctx context.Context, {{.in}}) (int64, error) {

    {{if .withCache}}
        data,err := m.FindOneBy{{.upperField}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
        return 0, err
        }
        {{.keys}}
        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            delTx := db.Where("{{.pkNameWrap}}", {{.lowerStartCamelField}}).Delete(&{{.upperStartCamelObject}}{})
            return delTx.RowsAffected,delTx.Error
        },{{.keyNames}})

    {{else}}
    {{/*var resp {{.upperStartCamelObject}}*/}}
    {{/*query := fmt.Sprintf("select %s from %s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)*/}}
    {{/*err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelField}})*/}}
    {{/*switch err {*/}}
    {{/*case nil:*/}}
    {{/*    return &resp, nil*/}}
    {{/*case sqlx.ErrNotFound:*/}}
    {{/*    return nil, ErrNotFound*/}}
    {{/*default:*/}}
    {{/*    return nil, err*/}}
    {{/*}*/}}
    {{end}}
    }