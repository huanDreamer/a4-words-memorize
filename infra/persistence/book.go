package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"words/domain/repository"
)

type MBook struct {
	// Id   string `bson:"_id"`
	BookId string `bson:"bookId"`
	Name   string `bson:"name"`
}

func (m MBook) CreateBook(ctx context.Context) (err error) {
	_, err = repository.GetCollection("book").InsertOne(ctx, m)
	return err
}

func (m MBook) BookInfo(ctx context.Context, bookId string) (result MBook, err error) {
	r := repository.GetCollection("book").FindOne(ctx, bson.D{{"bookId", bookId}})
	err = r.Decode(&result)
	return
}

func (m MBook) List(ctx context.Context, bookId string) (result []MBook, err error) {
	var filter interface{}
	if bookId != "" {
		filter = bson.D{{"bookId", bookId}}
	} else {
		filter = bson.D{}
	}
	cursor, err := repository.GetCollection("book").Find(ctx, filter)
	if err != nil {
		return
	}
	result = make([]MBook, 0)
	err = cursor.All(ctx, &result)
	return result, err
}
