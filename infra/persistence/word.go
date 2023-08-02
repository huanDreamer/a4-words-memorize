package persistence

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"words/domain/repository"
)

type MWord struct {
	WordRank int    `bson:"wordRank"`
	HeadWord string `bson:"headWord"`
	Content  struct {
		Word struct {
			WordHead string `bson:"wordHead"`
			WordId   string `bson:"wordId"`
			Content  struct {
				Sentence struct {
					Sentences []struct {
						SContent string `bson:"sContent"`
						SCn      string `bson:"sCn"`
					} `bson:"sentences"`
					Desc string `bson:"desc"`
				} `bson:"sentence"`
				Usphone string `bson:"usphone"`
				Syno    struct {
					Synos []struct {
						Pos  string `bson:"pos"`
						Tran string `bson:"tran"`
						Hwds []struct {
							W string `bson:"w"`
						} `bson:"hwds"`
					} `bson:"synos"`
					Desc string `bson:"desc"`
				} `bson:"syno"`
				Ukphone  string `bson:"ukphone"`
				Ukspeech string `bson:"ukspeech"`
				Phrase   struct {
					Phrases []struct {
						PContent string `bson:"pContent"`
						PCn      string `bson:"pCn"`
					} `bson:"phrases"`
					Desc string `bson:"desc"`
				} `bson:"phrase"`
				RelWord struct {
					Rels []struct {
						Pos   string `bson:"pos"`
						Words []struct {
							Hwd  string `bson:"hwd"`
							Tran string `bson:"tran"`
						} `bson:"words"`
					} `bson:"rels"`
					Desc string `bson:"desc"`
				} `bson:"relWord"`
				Usspeech string `bson:"usspeech"`
				Trans    []struct {
					TranCn    string `bson:"tranCn"`
					DescOther string `bson:"descOther"`
					Pos       string `bson:"pos"`
					DescCn    string `bson:"descCn"`
					TranOther string `bson:"tranOther"`
				} `bson:"trans"`
			} `bson:"content"`
		} `bson:"word"`
	} `bson:"content"`
	BookId string `bson:"bookId"`
}

func (m MWord) CreateWords(ctx context.Context, ws []MWord) (err error) {
	d := make([]interface{}, len(ws))
	for i, w := range ws {
		d[i] = w
	}
	_, err = repository.GetCollection("words").InsertMany(ctx, d)
	return
}

func (m MWord) WordNums(ctx context.Context, bookId string) (nums map[string]int, err error) {
	// 定义 Group By 和计算总数的聚合管道
	pipeline := bson.A{
		bson.D{{"$group", bson.D{{"_id", "$bookId"}, {"count", bson.D{{"$sum", 1}}}}}},
	}

	if bookId != "" {
		pipeline = append(pipeline, bson.D{{"$match", bson.D{{"bookId", bookId}}}})
	}

	// 执行聚合操作
	cursor, err := repository.GetCollection("words").Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Error executing aggregation:", err)
		return
	}
	// 迭代结果并输出
	defer cursor.Close(ctx)

	nums = make(map[string]int)

	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			fmt.Println("Error decoding result:", err)
			return
		}

		id := result["_id"]
		count := result["count"]
		nums[cast.ToString(id)] = cast.ToInt(count)
	}

	if err = cursor.Err(); err != nil {
		fmt.Println("Error iterating cursor:", err)
		return
	}
	return nums, nil
}

func (m MWord) FindByBook(ctx context.Context, bookId string, excludes []string, includes []string, num int64) (result []MWord, err error) {
	filter := bson.M{
		"bookId": bookId,
	}

	if len(excludes) > 0 {
		filter["headWord"] = bson.M{"$nin": excludes}
	}
	if len(includes) > 0 {
		filter["headWord"] = bson.M{"$in": includes}
	}

	// 定义查询选项，包含 limit
	findOptions := options.Find()
	findOptions.SetLimit(num)

	cursor, err := repository.GetCollection("words").Find(ctx, filter, findOptions)
	if err != nil {
		return
	}
	result = make([]MWord, 0)
	err = cursor.All(ctx, &result)
	return result, err
}
