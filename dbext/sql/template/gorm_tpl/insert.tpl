func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}},delCacheKeys ...string) (int64, error) {
{{if .withCache}}{{.keys}}
		afterDel := true
		cacheKeys := make([]string, 0, {{.keysLen}})
		cacheKeys = []string{ {{.uniqueCacheKeys}}}
		if len(delCacheKeys) > 0 {
			cacheKeys = append(cacheKeys, delCacheKeys...)
		}


		{{ if eq .pkDataType "string" }}
			if data.Id != "" {
				afterDel = false
				{{.pkCacheKeyExpression}}
				cacheKeys = append(cacheKeys, {{.pkCacheKey}})
			}
		{{else}}
			if data.Id != 0 {
				afterDel = false
				{{.pkCacheKeyExpression}}
				cacheKeys = append(cacheKeys, {{.pkCacheKey}})
			}
		{{end}}


		result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
			res := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Create(data)
			return res.RowsAffected, res.Error
		}, cacheKeys...)

		if err != nil {
			return 0, err
		}

		if afterDel {
			{{.pkCacheKeyExpression}}
			err = m.db.DelCacheKeys(ctx, {{.pkCacheKey}})
			if err != nil {
				return 0, err
			}
		}
		return result, err
{{else}}
		return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
			res := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Create(data)
			return res.RowsAffected, res.Error
		})
{{end}}
}
