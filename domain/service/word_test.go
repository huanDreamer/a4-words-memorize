package service

import (
	"bufio"
	"context"
	jsoniter "github.com/json-iterator/go"
	"os"
	"testing"
	"words/domain/entity"
)

func TestCreateWord(t *testing.T) {

	file, err := os.Open("../../source/IELTSluan_2.json")
	if err != nil {
		t.Error(err)
		return
	}
	// defer file.Close()

	words := make([]entity.Word, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		w := entity.Word{}
		err = jsoniter.UnmarshalFromString(line, &w)
		if err != nil {
			return
		}
		words = append(words, w)
	}

	if err = scanner.Err(); err != nil {
		t.Error(err)
	}

	err = NewWordsService(context.Background()).CreateWords(words)
	if err != nil {
		t.Error(err)
	}
}

func TestWordNums(t *testing.T) {
	nums, err := NewWordsService(context.Background()).WordNums()
	if err != nil {
		return
	}
	t.Log(nums)
}
