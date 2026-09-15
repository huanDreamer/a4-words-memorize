// Command import 把 source/ 目录下的一行一词（JSONL）词典导入本地文件存储，
// 并可选地生成一批「最近学习记录」作为测试数据。
//
// 用法：
//
//	go run ./cmd/import                     # 导入 source 下全部词典到 ./data
//	go run ./cmd/import -data ./data -plans # 同时生成示例学习计划
//	go run ./cmd/import -sample 200         # 每本书只导入前 200 个单词（快速造小数据集）
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"words/application"
	"words/domain/entity"
	"words/domain/repository"
	"words/domain/service"
)

// bookMeta 是每本词典对应的单词书元信息（书 id -> 展示名）。
var bookMeta = map[string]string{
	"CET4luan_1":  "四级词汇",
	"IELTSluan_2": "雅思词汇",
}

func main() {
	var (
		dataDir   = flag.String("data", "data", "数据目录")
		sourceDir = flag.String("source", "source", "词典源目录（*.json，一行一个单词）")
		sample    = flag.Int("sample", 0, "每本书最多导入多少个单词，0 表示全部")
		withPlans = flag.Bool("plans", false, "是否额外生成示例学习计划")
		planCount = flag.Int("plan-count", 5, "每本书生成的示例学习计划数量")
	)
	flag.Parse()

	repository.Default = repository.NewStore(*dataDir)
	ctx := context.Background()

	files, err := filepath.Glob(filepath.Join(*sourceDir, "*.json"))
	if err != nil || len(files) == 0 {
		fmt.Fprintf(os.Stderr, "在 %s 下没有找到词典文件: %v\n", *sourceDir, err)
		os.Exit(1)
	}
	sort.Strings(files)

	totalImported := 0
	bookWordCount := map[string]int{}
	for _, file := range files {
		words, err := readWords(file, *sample)
		if err != nil {
			fmt.Fprintf(os.Stderr, "读取 %s 失败: %v\n", file, err)
			os.Exit(1)
		}
		if len(words) == 0 {
			continue
		}
		bookId := words[0].BookId

		// 1. 建单词书
		name := bookMeta[bookId]
		if name == "" {
			name = bookId
		}
		if err := service.NewBookService(ctx).CreateBook(entity.Book{BookId: bookId, Name: name}); err != nil {
			fmt.Fprintf(os.Stderr, "创建单词书 %s 失败: %v\n", bookId, err)
			os.Exit(1)
		}

		// 2. 导入单词（按批写入，避免一次性占用过多内存）
		const batch = 500
		for i := 0; i < len(words); i += batch {
			end := i + batch
			if end > len(words) {
				end = len(words)
			}
			if err := service.NewWordsService(ctx).CreateWords(words[i:end]); err != nil {
				fmt.Fprintf(os.Stderr, "导入 %s 第 %d 批失败: %v\n", bookId, i/batch, err)
				os.Exit(1)
			}
		}
		bookWordCount[bookId] = len(words)
		totalImported += len(words)
		fmt.Printf("已导入 %-14s %-6s %5d 个单词\n", bookId, name, len(words))
	}

	// 3. 可选：生成示例学习计划
	if *withPlans {
		for bookId := range bookWordCount {
			for i := 0; i < *planCount; i++ {
				if _, err := application.NewStudyApplication(ctx).GenerateStudyPlan(bookId); err != nil {
					fmt.Fprintf(os.Stderr, "为 %s 生成学习计划失败: %v\n", bookId, err)
					os.Exit(1)
				}
			}
			fmt.Printf("已为 %-14s 生成 %d 条学习计划\n", bookId, *planCount)
		}
	}

	fmt.Printf("\n完成：共导入 %d 个单词，数据目录 %s\n", totalImported, repository.Default.Dir())
}

// readWords 读取一行一个单词的 JSONL 文件，最多返回 limit 个（limit<=0 表示不限）。
func readWords(file string, limit int) ([]entity.Word, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var words []entity.Word
	scanner := bufio.NewScanner(f)
	// 单行可能很长，放大扫描缓冲
	scanner.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var w entity.Word
		if err := jsoniter.UnmarshalFromString(line, &w); err != nil {
			return nil, fmt.Errorf("解析第 %d 行失败: %w", len(words)+1, err)
		}
		if w.BookId == "" {
			w.BookId = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		}
		words = append(words, w)
		if limit > 0 && len(words) >= limit {
			break
		}
	}
	return words, scanner.Err()
}
