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

	// 排除已经学过的单词，并随机挑选 25 个，避免每次学习顺序都一样
	words, err := service.NewWordsService(app.ctx).FindByBookRandom(bookId, studied, 25, true)
	if err != nil {
		return
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
	transMap, phoneMap := wordsToMaps(wordDetails)
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
			WordTrans: transMap[w.HeadWord],
			Phone:     phoneMap[w.HeadWord],
			Rank:      w.Rank,
			Mark:      w.Mark,
		}
	}
	return result, nil
}

// wordsToMaps 一次遍历，同时产出「释义」与「音标」两张表。
func wordsToMaps(words []entity.Word) (trans map[string][]string, phone map[string]string) {
	trans = make(map[string][]string, len(words))
	phone = make(map[string]string, len(words))
	for _, w := range words {
		c := w.Content.Word.Content
		for _, tran := range c.Trans {
			trans[w.HeadWord] = append(trans[w.HeadWord], fmt.Sprintf("%s %s", tran.Pos, tran.TranCn))
		}
		// 音标优先美式，缺失时退回英式
		if p := c.Usphone; p != "" {
			phone[w.HeadWord] = "/" + p + "/"
		} else if p := c.Ukphone; p != "" {
			phone[w.HeadWord] = "/" + p + "/"
		}
	}
	return trans, phone
}

// 标记为已学完
func (app StudyApplication) MarkStudied(planId int64) (err error) {
	return service.NewStudyPlanService(app.ctx).UpdateStatus(planId, "已完成")
}
