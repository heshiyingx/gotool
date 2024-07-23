FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error)
DeleteOneBy{{.upperField}}(ctx context.Context, {{.in}},delCacheKeys ...string) (int64, error)
UpdateOneBy{{.upperField}}(ctx context.Context, {{.in}},updateObj *{{.upperStartCamelObject}},delCacheKeys []string,fields ...string) (int64, error)