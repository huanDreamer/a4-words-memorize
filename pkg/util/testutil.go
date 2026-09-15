package util

import (
	"os"

	"words/domain/repository"
)

// SetupTestStore 把全局默认存储切换到临时目录，返回恢复函数。
// 供各测试包的 TestMain 使用，避免测试污染真实的 data 目录。
func SetupTestStore() (cleanup func()) {
	dir, err := os.MkdirTemp("", "words-test-*")
	if err != nil {
		panic(err)
	}
	old := repository.Default
	repository.Default = repository.NewStore(dir)
	return func() {
		repository.Default = old
		_ = os.RemoveAll(dir)
	}
}
