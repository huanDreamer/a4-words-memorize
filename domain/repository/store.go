package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store 是一个基于本地 JSON 文件的极简数据存储。
// 每个集合对应数据目录下的一个 .json 文件，读时整体加载、写时整体覆盖。
//
// 并发模型：
//   - Load 是**无锁**的。写入采用「写临时文件 + 原子 rename」，因此读方
//     看到的永远是某个完整版本，不会读到半个文件，无需加锁。
//   - 需要「读-改-写」时，必须用 WithLock 把整段包起来，以串行化同一集合
//     的并发写，避免丢失更新。
//   - Save 本身无锁，只应在 WithLock 内部（或确定无并发时）调用。
type Store struct {
	dir string

	mu       sync.Mutex // 保护 urbs 自身的 map
	urbs     map[string]*sync.Mutex
	backup   bool
	muBackup sync.Mutex
}

const defaultDir = "data"

var (
	ErrNotFound = errors.New("record not found")
)

// NewStore 创建一个本地文件存储。dir 为空时使用默认 data 目录。
func NewStore(dir string) *Store {
	if strings.TrimSpace(dir) == "" {
		dir = defaultDir
	}
	s := &Store{
		dir:    filepath.Clean(dir),
		urbs:   make(map[string]*sync.Mutex),
		backup: true,
	}
	_ = os.MkdirAll(s.dir, 0o755)
	return s
}

// filePath 返回集合对应的文件路径。
func (s *Store) filePath(collection string) string {
	return filepath.Join(s.dir, collection+".json")
}

// lock 返回并缓存某个集合的独立互斥锁，避免不同集合相互阻塞。
func (s *Store) lock(collection string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.urbs[collection]; ok {
		return m
	}
	m := &sync.Mutex{}
	s.urbs[collection] = m
	return m
}

// Save 将 v 序列化后原子写入集合文件（先写临时文件再 rename）。
// 无锁，需要并发安全时请放进 WithLock。
func (s *Store) Save(collection string, v interface{}) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	path := s.filePath(collection)
	tmp := path + ".tmp" // 同名即可：写操作已被 WithLock 串行化

	s.muBackup.Lock()
	if s.backup {
		// 复制而非 rename，避免出现「文件短暂不存在」的窗口破坏无锁读
		if _, statErr := os.Stat(path); statErr == nil {
			_ = copyFile(path, path+".bak")
		}
	}
	s.muBackup.Unlock()

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Load 读取集合文件并反序列化到 v。文件不存在或为空时视为空集合，返回 nil。
// 无锁：依赖 Save 的原子 rename 保证读到完整内容。
func (s *Store) Load(collection string, v interface{}) error {
	data, err := os.ReadFile(s.filePath(collection))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// WithLock 在持有 collection 对应锁的情况下执行 fn，用于串行化「读-改-写」。
// 注意：fn 内部请直接调用 Load/Save，它们不会再取锁，因此不会自锁。
func (s *Store) WithLock(collection string, fn func() error) error {
	m := s.lock(collection)
	m.Lock()
	defer m.Unlock()
	return fn()
}

// Dir 返回当前数据目录。
func (s *Store) Dir() string { return s.dir }

// copyFile 复制文件内容，用于生成 .bak 备份。
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// 全局默认存储实例，供 persistence 层使用。可在启动/测试时替换。
var Default = NewStore("")
