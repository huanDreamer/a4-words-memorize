package entity

import "words/infra/persistence"

type Book struct {
	BookId string
	Name   string
}

func (b Book) ToModel() (mbook persistence.MBook) {
	return persistence.MBook{
		BookId: b.BookId,
		Name:   b.Name,
	}
}

func (b Book) FromModel(m persistence.MBook) Book {
	b.BookId = m.BookId
	b.Name = m.Name
	return b
}
