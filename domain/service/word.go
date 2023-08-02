package service

import (
	"context"
	"math"
	"words/domain/entity"
	"words/infra/persistence"
)

type WordService struct {
	ctx context.Context
}

func NewWordsService(ctx context.Context) *WordService {
	return &WordService{ctx: ctx}
}

func (svc WordService) CreateWords(w []entity.Word) error {
	d := make([]persistence.MWord, len(w))
	for i, ww := range w {
		d[i] = ww.MWord
	}
	return persistence.MWord{}.CreateWords(svc.ctx, d)
}

func (svc WordService) WordNums() (nums map[string]int, err error) {
	return persistence.MWord{}.WordNums(svc.ctx, "")
}

// 查找单词
func (svc WordService) FindByBook(bookId string, excludes []string, num int64) (result []entity.Word, err error) {
	words, err := persistence.MWord{}.FindByBook(svc.ctx, bookId, excludes, nil, num)
	if err != nil {
		return nil, err
	}
	for _, word := range words {
		result = append(result, entity.Word{}.FromModel(word))
	}
	return result, nil
}

// 查找单词
func (svc WordService) FindWordsDetail(bookId string, includeWords []string) (result []entity.Word, err error) {
	words, err := persistence.MWord{}.FindByBook(svc.ctx, bookId, nil, includeWords, math.MaxInt)
	if err != nil {
		return nil, err
	}
	for _, word := range words {
		result = append(result, entity.Word{}.FromModel(word))
	}
	return result, nil
}
