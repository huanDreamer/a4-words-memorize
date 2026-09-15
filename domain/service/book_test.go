package service

import (
	"context"
	"os"
	"testing"

	"words/domain/entity"
	"words/pkg/util"
)

func TestMain(m *testing.M) {
	cleanup := util.SetupTestStore()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestInitBook(t *testing.T) {
	b := entity.Book{
		BookId: "IELTSluan_2",
		Name:   "雅思词汇",
	}

	err := NewBookService(context.Background()).CreateBook(b)
	if err != nil {
		t.Error(err)
	}
}
