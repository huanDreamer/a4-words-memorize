package persistence

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"words/domain/repository"
)

type MStudyPlan struct {
	Date   string  `bson:"date"`
	PlanId int64   `bson:"planId"`
	BookId string  `bson:"bookId"`
	Num    int     `bson:"num"`
	Status string  `bson:"status"`
	Words  []VWord `bson:"words"`
}

type VWord struct {
	HeadWord string `bson:"headWord"`
	Rank     int    `bson:"rank"` // 顺序
	Mark     int    `bson:"mark"` // 标记
}

func (m MStudyPlan) Create(ctx context.Context) (planId int64, err error) {
	_, err = repository.GetCollection("study_plan").InsertOne(ctx, m)
	return m.PlanId, err
}

func (m MStudyPlan) FindPlan(ctx context.Context, planId int64) (result MStudyPlan, err error) {
	r := repository.GetCollection("study_plan").FindOne(ctx, bson.D{{"planId", planId}})
	err = r.Decode(&result)
	return
}

func (m MStudyPlan) Update(ctx context.Context) (err error) {
	update := bson.M{"$set": bson.M{"status": m.Status}}
	_, err = repository.GetCollection("study_plan").UpdateOne(ctx, bson.D{{"planId", m.PlanId}}, update)
	return
}

func (m MStudyPlan) FindByBook(ctx context.Context, bookId string, num int64) (result []MStudyPlan, err error) {
	var filter interface{}
	if bookId != "" {
		filter = bson.D{{"bookId", bookId}}
	} else {
		filter = bson.D{}
	}
	// 定义查询选项，包含 limit
	findOptions := options.Find()
	findOptions.SetLimit(num).SetSort(bson.D{{"planId", 1}}) // -1 降序 1 升序

	cursor, err := repository.GetCollection("study_plan").Find(ctx, filter)
	if err != nil {
		return
	}
	result = make([]MStudyPlan, 0)
	err = cursor.All(ctx, &result)
	return result, err
}
