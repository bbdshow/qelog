package mongoclient

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// uri 可以包含连接的所有信息 包括设置连接读写超时等
func NewMongoClientByURI(ctx context.Context, uri string) (*mongo.Client, error) {
	opt := options.Client().ApplyURI(uri)
	return NewMongoClient(ctx, opt)
}

func NewMongoClient(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// ping
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Ping(ctxTimeout, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type Database struct {
	*mongo.Database
}

func (db *Database) FindOne(ctx context.Context, coll *mongo.Collection, filter interface{}, ret interface{}, opt ...*options.FindOneOptions) (bool, error) {
	sRet := coll.FindOne(ctx, filter, opt...)
	if sRet.Err() != nil {
		if sRet.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, sRet.Err()
	}

	if err := sRet.Decode(ret); err != nil {
		return false, err
	}
	return true, nil
}

func (db *Database) Find(ctx context.Context, coll *mongo.Collection, filter bson.M, records interface{}, opt ...*options.FindOptions) error {
	cursor, err := coll.Find(ctx, filter, opt...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, records); err != nil {
		return err
	}
	return nil
}

func (db *Database) FindAndCount(ctx context.Context, coll *mongo.Collection, filter interface{}, result interface{}, opt ...*options.FindOptions) (int64, error) {
	c, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	cursor, err := coll.Find(ctx, filter, opt...)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, result); err != nil {
		return 0, err
	}
	return c, nil
}

type Index struct {
	Collection         string
	Name               string // 指定索引名称
	Keys               bson.M
	Unique             bool  // 唯一索引
	Background         bool  // 非阻塞创建索引
	ExpireAfterSeconds int32 // 多少秒后过期
}

func (i Index) Validate() error {
	if i.Collection == "" {
		return errors.New("collection required")
	}
	if len(i.Keys) == 0 {
		return errors.New("keys required")
	}
	return nil
}

func UpsertCollectionIndexMany(database *mongo.Database, indexMany []Index) error {
	indexModels := make(map[string][]mongo.IndexModel)
	for _, index := range indexMany {
		if err := index.Validate(); err != nil {
			return err
		}
		model := mongo.IndexModel{
			Keys: index.Keys,
		}
		opt := options.Index()
		if index.Name != "" {
			opt.SetName(index.Name)
		}
		opt.SetUnique(index.Unique)
		opt.SetBackground(index.Background)

		if index.ExpireAfterSeconds > 0 {
			opt.SetExpireAfterSeconds(index.ExpireAfterSeconds)
		}

		model.Options = opt

		v, ok := indexModels[index.Collection]
		if ok {
			indexModels[index.Collection] = append(v, model)
		} else {
			indexModels[index.Collection] = []mongo.IndexModel{model}
		}
	}

	for collection, index := range indexModels {
		_, err := database.Collection(collection).Indexes().CreateMany(context.Background(), index)
		if err != nil {
			return err
		}
	}

	return nil
}
