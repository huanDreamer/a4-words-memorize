package service

import (
	"context"
	"testing"

	"words/domain/entity"
)

func TestGeneratePlanExcludesStudied(t *testing.T) {
	ctx := context.Background()
	bookId := "test_book"

	// 准备 30 个单词
	var ws []entity.Word
	for i := 0; i < 30; i++ {
		w := entity.Word{}
		w.BookId = bookId
		w.HeadWord = "w" + string(rune('a'+i%26)) + string(rune('0'+i/26))
		w.WordRank = i + 1
		ws = append(ws, w)
	}
	if err := NewWordsService(ctx).CreateWords(ws); err != nil {
		t.Fatalf("CreateWords: %v", err)
	}

	// 第一次生成计划：应取到 25 个
	svc := NewStudyPlanService(ctx)
	planId, err := svc.GeneratePlan(entity.StudyPlan{
		Date:   "2024-01-01",
		BookId: bookId,
		Status: "未完成",
		Words:  ws[:25],
	})
	if err != nil {
		t.Fatalf("GeneratePlan: %v", err)
	}
	plan, err := svc.FindPlan(planId)
	if err != nil {
		t.Fatalf("FindPlan: %v", err)
	}
	if plan.Num != 25 {
		t.Fatalf("want 25 words, got %d", plan.Num)
	}

	// 已学单词应能从查询中排除
	studied, err := svc.GetStudiedWords(bookId)
	if err != nil {
		t.Fatalf("GetStudiedWords: %v", err)
	}
	if len(studied) != 25 {
		t.Fatalf("want 25 studied words, got %d", len(studied))
	}

	rest, err := NewWordsService(ctx).FindByBook(bookId, studied, 25)
	if err != nil {
		t.Fatalf("FindByBook: %v", err)
	}
	if len(rest) == 0 {
		t.Fatal("expected remaining unstudied words")
	}
	for _, w := range rest {
		for _, s := range studied {
			if w.HeadWord == s {
				t.Fatalf("studied word %q should have been excluded", w.HeadWord)
			}
		}
	}
}

func TestFindByBookShuffleVaries(t *testing.T) {
	ctx := context.Background()
	bookId := "shuffle_book"
	var ws []entity.Word
	for i := 0; i < 100; i++ {
		w := entity.Word{}
		w.BookId = bookId
		w.HeadWord = "word" + string(rune('A'+i%26)) + string(rune('a'+i/26))
		ws = append(ws, w)
	}
	if err := NewWordsService(ctx).CreateWords(ws); err != nil {
		t.Fatal(err)
	}

	first, err := NewWordsService(ctx).FindByBookRandom(bookId, nil, 25, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 25 {
		t.Fatalf("want 25, got %d", len(first))
	}

	// 多取几次，至少有一次顺序不同（极低概率全相同）
	same := true
	for i := 0; i < 5 && same; i++ {
		next, _ := NewWordsService(ctx).FindByBookRandom(bookId, nil, 25, true)
		for j := range next {
			if next[j].HeadWord != first[j].HeadWord {
				same = false
				break
			}
		}
	}
	if same {
		t.Fatal("shuffle produced identical order 6 times; randomization likely broken")
	}
}

func TestMarkStudiedUpdatesStatus(t *testing.T) {
	ctx := context.Background()
	svc := NewStudyPlanService(ctx)
	w := entity.Word{}
	w.BookId = "b"
	w.HeadWord = "hello"
	planId, err := svc.GeneratePlan(entity.StudyPlan{
		Date: "2024-01-01", BookId: "b", Status: "未完成", Words: []entity.Word{w},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.UpdateStatus(planId, "已完成"); err != nil {
		t.Fatalf("UpdateStatus: %v", err)
	}
	plan, err := svc.FindPlan(planId)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != "已完成" {
		t.Fatalf("want 已完成, got %q", plan.Status)
	}
}

func TestFindRecentPlansNewestFirst(t *testing.T) {
	ctx := context.Background()
	svc := NewStudyPlanService(ctx)
	var lastId int64
	for i := 0; i < 3; i++ {
		id, err := svc.GeneratePlan(entity.StudyPlan{
			Date: "2024-01-01", BookId: "bb", Status: "未完成",
			Words: []entity.Word{{}},
		})
		if err != nil {
			t.Fatal(err)
		}
		lastId = id
	}
	plans, err := svc.FindPlans("bb", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 3 {
		t.Fatalf("want 3 plans, got %d", len(plans))
	}
	if plans[0].PlanId != lastId {
		t.Fatalf("newest plan should be first: want %d, got %d", lastId, plans[0].PlanId)
	}
	for i := 1; i < len(plans); i++ {
		if plans[i-1].PlanId < plans[i].PlanId {
			t.Fatal("plans not sorted newest-first")
		}
	}
}
