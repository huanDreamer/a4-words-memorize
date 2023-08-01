package entity

import (
	"time"
	"words/infrastructure/persistence"
)

type StudyPlan struct {
	Date   string
	BookId string
	Status string
	Words  []Word
}

func (s StudyPlan) ToModel() (m persistence.MStudyPlan) {
	m.Date = s.Date
	m.PlanId = time.Now().UnixMilli()
	m.BookId = s.BookId
	m.Num = len(s.Words)
	m.Status = s.Status

	m.Words = make([]persistence.VWord, len(s.Words))
	for i, w := range s.Words {
		w.Rank = i
		m.Words[i] = persistence.VWord{
			HeadWord: w.HeadWord,
			Rank:     w.Rank,
			Mark:     0,
		}
	}
	return m
}

type StudyPlanInfo struct {
	persistence.MStudyPlan
}

func (s StudyPlanInfo) FromModel(m persistence.MStudyPlan) StudyPlanInfo {
	s.MStudyPlan = m
	return s
}
