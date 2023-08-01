package application

import (
	"context"
	"time"
	"words/domain/service"
	"words/interfaces/response"
)

type BookApplication struct {
	ctx context.Context
}

func NewBookApplication(ctx context.Context) *BookApplication {
	return &BookApplication{ctx: ctx}
}

// 列表
func (app BookApplication) BookList() (result []response.BookInfo) {

	// 书籍信息
	books, err := service.NewBookService(app.ctx).List()
	if err != nil {
		return
	}
	// 单词数量
	wordNums, err := service.NewWordsService(app.ctx).WordNums()
	if err != nil {
		return nil
	}

	// 最近的学习计划
	plans, err := service.NewStudyPlanService(app.ctx).FindPlans("", 30)
	if err != nil {
		return nil
	}
	planMap := make(map[string][]response.StudyPlan)
	for _, plan := range plans {
		planMap[plan.BookId] = append(planMap[plan.BookId], response.StudyPlan{
			Time:   time.UnixMilli(plan.PlanId).Format("2006-01-02 15:04:05"),
			PlanId: plan.PlanId,
			Status: plan.Status,
			Num:    plan.Num,
		})
	}

	// 组装
	for _, b := range books {
		result = append(result, response.BookInfo{
			BookId:     b.BookId,
			Name:       b.Name,
			WordNum:    wordNums[b.BookId],
			StudyPlans: planMap[b.BookId],
		})
	}

	return
}
