package application

import (
	"context"
	"log"
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

// IndexData 首页所需数据：单词本列表（含各自最近的学习记录）与顶部统计。
func (app BookApplication) IndexData() (result response.IndexData) {
	result.Books = app.BookList()
	result.Now = time.Now().Format("2006-01-02 15:04")
	result.TotalBooks = len(result.Books)
	for _, b := range result.Books {
		result.TotalWords += b.WordNum
		result.TotalPlans += len(b.StudyPlans)
	}
	return result
}

// BookList 列表
func (app BookApplication) BookList() (result []response.BookInfo) {

	// 书籍信息
	books, err := service.NewBookService(app.ctx).List()
	if err != nil {
		log.Printf("读取单词本失败: %v", err)
		return nil
	}

	// 单词数量
	wordNums, err := service.NewWordsService(app.ctx).WordNums()
	if err != nil {
		// 统计失败不该让整页变空，降级为 0 继续渲染
		log.Printf("统计单词数量失败: %v", err)
	}

	// 最近的学习计划
	plans, err := service.NewStudyPlanService(app.ctx).FindPlans("", 30)
	if err != nil {
		log.Printf("读取最近学习记录失败: %v", err)
	}
	planMap := make(map[string][]response.StudyPlan, len(books))
	for _, plan := range plans {
		planMap[plan.BookId] = append(planMap[plan.BookId], response.StudyPlan{
			Time:   time.UnixMilli(plan.PlanId).Format("2006-01-02 15:04:05"),
			PlanId: plan.PlanId,
			Status: plan.Status,
			Num:    plan.Num,
		})
	}

	// 组装
	result = make([]response.BookInfo, 0, len(books))
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
