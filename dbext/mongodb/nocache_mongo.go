package gormdb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

type (
	MongoDB struct {
		singleFlight *singleflight.Group
		db           *mongo.Client
		collection   *mongo.Collection
	}
)

func MustNewMongoDB(c Config, opts ...Option) *MongoDB {
	mongoDB, err := NewMongoDB(c, opts...)
	if err != nil {
		log.Fatalf("NewCacheGormDB err:%v", err)
		return nil
	}
	return mongoDB
}

func NewMongoDB(c Config, opts ...Option) (*MongoDB, error) {

	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(c.URI).
		//ApplyURI("mongodb://localhost:27017").
		SetMaxPoolSize(150).
		SetMinPoolSize(10).
		SetConnectTimeout(time.Second * 10).
		SetRetryReads(true).
		SetRetryWrites(true)
	if c.User != "" && c.Pwd != "" {
		clientOptions.SetAuth(options.Credential{
			AuthSource:  "admin",
			Username:    c.User,
			Password:    c.Pwd,
			PasswordSet: true,
		})
	}
	for _, opt := range opts {
		opt(clientOptions)
	}

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Cannot connect to MongoDB!")
	}
	collection := client.Database(c.DataBase).Collection(c.Collection)
	return &MongoDB{
		singleFlight: &singleflight.Group{},
		db:           client,
		collection:   collection,
	}, nil
}

func (cg *MongoDB) QueryCtx(ctx context.Context, result any, key string, fn QueryCtxFn) error {
	return cg.takeCtx(ctx, key, result, fn)
}
func (cg *MongoDB) QueryNoCacheCtx(ctx context.Context, result any, fn QueryCtxFn) error {
	return fn(ctx, result, cg.collection)
}
func (cg *MongoDB) takeCtx(ctx context.Context, key string, result any, query QueryCtxFn) error {

	_, err, _ := cg.singleFlight.Do(key, func() (interface{}, error) {

		if err := query(ctx, result, cg.collection); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, mongo.ErrNoDocuments
			} else {
				return nil, err
			}
		}
		return result, nil
	})
	return err
}
func (cg *MongoDB) ExecCtx(ctx context.Context, execFn ExecCtxFn) (int64, error) {

	result, err := execFn(ctx, cg.collection)
	if err != nil {
		return 0, err
	}
	return result, nil
}
