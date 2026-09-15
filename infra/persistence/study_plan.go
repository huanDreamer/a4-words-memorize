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

// Create 追加一条学习计划并返回最终使用的 planId（唯一且单调递增）。
//
// PlanId 由调用方按毫秒时间戳生成，毫秒内连续创建会撞号；早期实现靠 Sleep
// 规避。这里在锁内重新分配：若计划 id 不大于已存在的最大值，则顺延一位，
// 从而保证唯一、递增，也顺带获得稳定排序。
func (m MStudyPlan) Create(ctx context.Context) (planId int64, err error) {
	err = repository.Default.WithLock(studyPlanCollection, func() error {
		var plans []MStudyPlan
		if err := repository.Default.Load(studyPlanCollection, &plans); err != nil {
			return err
		}
		var maxId int64
		for _, p := range plans {
			if p.PlanId > maxId {
				maxId = p.PlanId
			}
		}
		if m.PlanId <= maxId {
			m.PlanId = maxId + 1
		}
		planId = m.PlanId
		plans = append(plans, m)
		return repository.Default.Save(studyPlanCollection, &plans)
	})
	return planId, err
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
// num<=0 表示不限；PlanId 基于毫秒时间戳，天然反映创建先后。
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

// FindRecentByBook 查询「最近的学习记录」：返回最新的 num 条计划，按 PlanId 降序（新→旧）。
func (m MStudyPlan) FindRecentByBook(ctx context.Context, bookId string, num int64) (result []MStudyPlan, err error) {
	plans, err := m.FindByBook(ctx, bookId, 0)
	if err != nil {
		return nil, err
	}
	// 反转成降序，从最新开始取
	result = make([]MStudyPlan, 0, len(plans))
	for i := len(plans) - 1; i >= 0; i-- {
		p := plans[i]
		if bookId != "" && p.BookId != bookId {
			continue
		}
		result = append(result, p)
		if num > 0 && int64(len(result)) >= num {
			break
		}
	}
	return result, nil
}
