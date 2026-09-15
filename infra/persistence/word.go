package persistence

import (
	"context"
	"math/rand"

	"words/domain/repository"
)

// MWord 是单词在文件存储中的模型。
// 相比原始词典，仅保留应用实际需要的字段（发音音标、中文释义&词性），
// 大幅缩小数据体积，也便于文件读写。
type MWord struct {
	WordRank int    `json:"wordRank"`
	HeadWord string `json:"headWord"`
	BookId   string `json:"bookId"`
	Content  struct {
		Word struct {
			WordHead string `json:"wordHead"`
			WordId   string `json:"wordId"`
			Content  struct {
				Usphone string `json:"usphone"`
				Ukphone string `json:"ukphone"`
				Trans   []struct {
					TranCn string `json:"tranCn"`
					Pos    string `json:"pos"`
				} `json:"trans"`
			} `json:"content"`
		} `json:"word"`
	} `json:"content"`
}

const wordCollection = "words"

// CreateWords 批量写入单词（按 bookId+headWord 幂等去重）。
func (m MWord) CreateWords(ctx context.Context, ws []MWord) (err error) {
	return repository.Default.WithLock(wordCollection, func() error {
		var all []MWord
		if err := repository.Default.Load(wordCollection, &all); err != nil {
			return err
		}
		seen := make(map[string]bool, len(all))
		for _, w := range all {
			seen[w.BookId+"|"+w.HeadWord] = true
		}
		for _, w := range ws {
			if seen[w.BookId+"|"+w.HeadWord] {
				continue
			}
			all = append(all, w)
			seen[w.BookId+"|"+w.HeadWord] = true
		}
		return repository.Default.Save(wordCollection, &all)
	})
}

// WordNums 统计每个单词本的单词数量。bookId 非空时只统计对应单词本。
func (m MWord) WordNums(ctx context.Context, bookId string) (nums map[string]int, err error) {
	var all []MWord
	if err = repository.Default.Load(wordCollection, &all); err != nil {
		return nil, err
	}
	nums = make(map[string]int)
	for _, w := range all {
		if bookId != "" && w.BookId != bookId {
			continue
		}
		nums[w.BookId]++
	}
	return nums, nil
}

// FindByBook 按单词本查找单词。
// excludes: 命中的 headWord 将被排除；includes: 仅返回命中的 headWord。
// num<=0 时返回全部；否则最多返回 num 个。
// shuffle 为 true 时，在满足条件的单词中随机挑选 num 个（用于打乱学习顺序）。
func (m MWord) FindByBook(ctx context.Context, bookId string, excludes []string, includes []string, num int64, shuffle bool) (result []MWord, err error) {
	var all []MWord
	if err = repository.Default.Load(wordCollection, &all); err != nil {
		return nil, err
	}
	excludeSet := toSet(excludes)
	includeSet := toSet(includes)
	hasExclude := len(excludes) > 0
	hasInclude := len(includes) > 0

	matched := make([]MWord, 0, len(all))
	for _, w := range all {
		if w.BookId != bookId {
			continue
		}
		if hasExclude && excludeSet[w.HeadWord] {
			continue
		}
		if hasInclude && !includeSet[w.HeadWord] {
			continue
		}
		matched = append(matched, w)
	}

	// 需要随机挑选且确实存在可挑选的余量时，先整体打乱再截断
	if shuffle && num > 0 && int64(len(matched)) > num {
		rand.Shuffle(len(matched), func(i, j int) { matched[i], matched[j] = matched[j], matched[i] })
	}

	if num > 0 && int64(len(matched)) > num {
		matched = matched[:num]
	}
	return matched, nil
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
