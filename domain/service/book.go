package service

import (
	"context"
	"words/domain/entity"
	"words/infrastructure/persistence"
)

type BookService struct {
	ctx context.Context
}

func NewBookService(ctx context.Context) *BookService {
	return &BookService{ctx: ctx}
}

func (svc BookService) CreateBook(b entity.Book) error {
	return b.ToModel().CreateBook(svc.ctx)
}

func (svc BookService) List() (result []entity.Book, err error) {
	books, err := persistence.MBook{}.List(svc.ctx, "")
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		result = append(result, entity.Book{}.FromModel(b))
	}
	return result, nil
}

func (svc BookService) BookInfo(bookId string) (b entity.Book, err error) {
	m, err := persistence.MBook{}.BookInfo(svc.ctx, bookId)
	if err != nil {
		return
	}
	return b.FromModel(m), err
}
