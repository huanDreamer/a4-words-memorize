package entity

import "words/infrastructure/persistence"

type Word struct {
	persistence.MWord
	Rank int
}

func (w Word) FromModel(m persistence.MWord) Word {
	w.MWord = m
	return w
}
