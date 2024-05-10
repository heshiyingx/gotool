package cacheext

import "context"

type (
	Cache interface {
		DelCtx(ctx context.Context, keys ...string) error
	}
)
