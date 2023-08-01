package service

import (
	"context"
	"errors"
	"math"
	"words/domain/entity"
	"words/infrastructure/persistence"
)

type StudyPlanService struct {
	ctx context.Context
}

func NewStudyPlanService(ctx context.Context) *StudyPlanService {
	return &StudyPlanService{ctx: ctx}
}

func (svc StudyPlanService) GetStudiedWords(bookId string) (result []string, err error) {
	plans, err := persistence.MStudyPlan{}.FindByBook(svc.ctx, bookId, math.MaxInt64)
	if err != nil {
		return nil, err
	}
	for _, plan := range plans {
		for _, word := range plan.Words {
			result = append(result, word.HeadWord)
		}
	}
	return result, nil
}

func (svc StudyPlanService) GeneratePlan(studyPlan entity.StudyPlan) (planId int64, err error) {
	m := studyPlan.ToModel()
	return m.Create(svc.ctx)
}

func (svc StudyPlanService) FindPlan(planId int64) (result entity.StudyPlanInfo, err error) {
	plan, err := persistence.MStudyPlan{}.FindPlan(svc.ctx, planId)
	if err != nil {
		return result, err
	}
	return result.FromModel(plan), err
}

func (svc StudyPlanService) UpdateStatus(planId int64, status string) (err error) {
	plan, err := persistence.MStudyPlan{}.FindPlan(svc.ctx, planId)
	if err != nil {
		return err
	}
	if plan.PlanId == 0 {
		return errors.New("计划未找到")
	}
	plan.Status = status
	return plan.Update(svc.ctx)
}

// 查找最近的学习计划
func (svc StudyPlanService) FindPlans(bookId string, num int64) (result []entity.StudyPlanInfo, err error) {
	plans, err := persistence.MStudyPlan{}.FindByBook(svc.ctx, bookId, num)
	if err != nil {
		return result, err
	}
	result = make([]entity.StudyPlanInfo, len(plans))
	for i, plan := range plans {
		result[i] = result[i].FromModel(plan)
	}
	return result, nil
}
