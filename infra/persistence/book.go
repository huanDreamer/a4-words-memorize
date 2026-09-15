package persistence

import (
	"context"
	"words/domain/repository"
)

type MBook struct {
	BookId string `json:"bookId"`
	Name   string `json:"name"`
}

const bookCollection = "book"

func (m MBook) CreateBook(ctx context.Context) (err error) {
	return repository.Default.WithLock(bookCollection, func() error {
		var books []MBook
		if err := repository.Default.Load(bookCollection, &books); err != nil {
			return err
		}
		// 幂等：已存在的书不重复写入
		for _, b := range books {
			if b.BookId == m.BookId {
				return nil
			}
		}
		books = append(books, m)
		return repository.Default.Save(bookCollection, &books)
	})
}

func (m MBook) BookInfo(ctx context.Context, bookId string) (result MBook, err error) {
	var books []MBook
	if err = repository.Default.Load(bookCollection, &books); err != nil {
		return
	}
	for _, b := range books {
		if b.BookId == bookId {
			return b, nil
		}
	}
	return result, repository.ErrNotFound
}

func (m MBook) List(ctx context.Context, bookId string) (result []MBook, err error) {
	if err = repository.Default.Load(bookCollection, &result); err != nil {
		return nil, err
	}
	if bookId != "" {
		filtered := make([]MBook, 0)
		for _, b := range result {
			if b.BookId == bookId {
				filtered = append(filtered, b)
			}
		}
		result = filtered
	}
	return result, nil
}
