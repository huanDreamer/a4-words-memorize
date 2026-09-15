package persistence

import (
	"context"
	"sort"

	"words/domain/repository"
)

type MStudyPlan struct {
	Date   string  `json:"date"`
	PlanId int64   `json:"planId"`
	BookId string  `json:"bookId"`
	Num    int     `json:"num"`
	Status string  `json:"status"`
	Words  []VWord `json:"words"`
}

type VWord struct {
	HeadWord string `json:"headWord"`
	Rank     int    `json:"rank"` // 顺序
	Mark     int    `json:"mark"` // 标记
}

const studyPlanCollection = "study_plan"

func (m MStudyPlan) Create(ctx context.Context) (planId int64, err error) {
	err = repository.Default.WithLock(studyPlanCollection, func() error {
		var plans []MStudyPlan
		if err := repository.Default.Load(studyPlanCollection, &plans); err != nil {
			return err
		}
		plans = append(plans, m)
		return repository.Default.Save(studyPlanCollection, &plans)
	})
	return m.PlanId, err
}

func (m MStudyPlan) FindPlan(ctx context.Context, planId int64) (result MStudyPlan, err error) {
	var plans []MStudyPlan
	if err = repository.Default.Load(studyPlanCollection, &plans); err != nil {
		return result, err
	}
	for _, p := range plans {
		if p.PlanId == planId {
			return p, nil
		}
	}
	return result, repository.ErrNotFound
}

func (m MStudyPlan) Update(ctx context.Context) (err error) {
	return repository.Default.WithLock(studyPlanCollection, func() error {
		var plans []MStudyPlan
		if err := repository.Default.Load(studyPlanCollection, &plans); err != nil {
			return err
		}
		for i := range plans {
			if plans[i].PlanId == m.PlanId {
				plans[i].Status = m.Status
				return repository.Default.Save(studyPlanCollection, &plans)
			}
		}
		return repository.ErrNotFound
	})
}

// FindByBook 查询指定单词本/全部的学习计划，按 PlanId 升序返回，最多 num 条。
func (m MStudyPlan) FindByBook(ctx context.Context, bookId string, num int64) (result []MStudyPlan, err error) {
	var plans []MStudyPlan
	if err = repository.Default.Load(studyPlanCollection, &plans); err != nil {
		return nil, err
	}
	sort.Slice(plans, func(i, j int) bool { return plans[i].PlanId < plans[j].PlanId })
	result = make([]MStudyPlan, 0, len(plans))
	count := int64(0)
	for _, p := range plans {
		if bookId != "" && p.BookId != bookId {
			continue
		}
		result = append(result, p)
		count++
		if num > 0 && count >= num {
			break
		}
	}
	return result, nil
}
