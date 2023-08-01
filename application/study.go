package application

import (
	"context"
	"fmt"
	"time"
	"words/domain/entity"
	"words/domain/service"
	"words/interfaces/response"
)

type StudyApplication struct {
	ctx context.Context
}

func NewStudyApplication(ctx context.Context) *StudyApplication {
	return &StudyApplication{ctx: ctx}
}

// 生成一个学习计划
func (app StudyApplication) GenerateStudyPlan(bookId string) (planId int64, err error) {
	// 查询已经学习过的
	studied, err := service.NewStudyPlanService(app.ctx).GetStudiedWords(bookId)
	if err != nil {
		return
	}

	// 排除已经学过的单词
	words, err := service.NewWordsService(app.ctx).FindByBook(bookId, studied, 25)
	if err != nil {
		return
	}
	wordMap := make(map[string][]string)
	for _, w := range words {
		for _, tran := range w.Content.Word.Content.Trans {
			wordMap[w.HeadWord] = append(wordMap[w.HeadWord], fmt.Sprintf("%s %s", tran.Pos, tran.TranCn))
		}
	}

	// 生成计划
	p := entity.StudyPlan{
		Date:   time.Now().Format("2006-01-02"),
		BookId: bookId,
		Words:  words,
		Status: "未完成",
	}
	planStudyService := service.NewStudyPlanService(app.ctx)
	return planStudyService.GeneratePlan(p)
}

// 获取一个学习计划
func (app StudyApplication) GetStudyPlan(planId int64) (result response.StudyWordResp, err error) {
	plan, err := service.NewStudyPlanService(app.ctx).FindPlan(planId)
	if err != nil {
		return
	}
	// 获取课本信息
	book, err := service.NewBookService(app.ctx).BookInfo(plan.BookId)
	if err != nil {
		return
	}
	// 获取单词信息
	words := make([]string, len(plan.Words))
	for i, w := range plan.Words {
		words[i] = w.HeadWord
	}
	wordDetails, err := service.NewWordsService(app.ctx).FindWordsDetail(plan.BookId, words)
	if err != nil {
		return
	}
	wordMap := wordsToMap(wordDetails)
	// 组装
	result = response.StudyWordResp{
		Date:   plan.Date,
		BookId: plan.BookId,
		PlanId: planId,
		Name:   book.Name,
		Num:    plan.Num,
		Status: plan.Status,
		Words:  make([]response.WordInfo, len(plan.Words)),
	}
	for i, w := range plan.Words {
		result.Words[i] = response.WordInfo{
			HeadWord:  w.HeadWord,
			WordTrans: wordMap[w.HeadWord],
			Rank:      w.Rank,
			Mark:      w.Mark,
		}
	}
	return result, nil
}

func wordsToMap(words []entity.Word) map[string][]string {
	wordMap := make(map[string][]string)
	for _, w := range words {
		for _, tran := range w.Content.Word.Content.Trans {
			wordMap[w.HeadWord] = append(wordMap[w.HeadWord], fmt.Sprintf("%s %s", tran.Pos, tran.TranCn))
		}
	}
	return wordMap
}

// 标记为已学完
func (app StudyApplication) MarkStudied(planId int64) (err error) {
	return service.NewStudyPlanService(app.ctx).UpdateStatus(planId, "已完成")
}
