package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var MongoClient *mongo.Client

func InitMongo() {
	MongoClient = newMongo()
}

func newMongo() *mongo.Client {
	opt := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	opt.Auth = &options.Credential{
		Username: "admin",
		Password: "123456",
	}
	opt.SetLocalThreshold(3 * time.Second)                  // 只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(time.Duration(30) * time.Second) // 指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(100)                                 // 使用最大的连接数
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}
	return client
}

func GetCollection(c string) *mongo.Collection {
	return MongoClient.Database("words").Collection(c)
}

// 关闭数据库连接
func CloseMongo() {
	_ = MongoClient.Disconnect(context.Background())
}
