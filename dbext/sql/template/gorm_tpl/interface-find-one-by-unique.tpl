FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error)
DeleteOneBy{{.upperField}}(ctx context.Context, {{.in}}) (int64, error)
UpdateOneBy{{.upperField}}(ctx context.Context, {{.in}},updateObj *{{.upperStartCamelObject}},fields ...string) (int64, error)