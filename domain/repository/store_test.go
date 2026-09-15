package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

type record struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	return NewStore(filepath.Join(dir, "data"))
}

func TestSaveLoadRoundTrip(t *testing.T) {
	s := newTestStore(t)
	in := []record{{1, "a"}, {2, "b"}}
	if err := s.Save("c", in); err != nil {
		t.Fatalf("Save: %v", err)
	}
	var out []record
	if err := s.Load("c", &out); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(out) != 2 || out[0].Name != "a" || out[1].ID != 2 {
		t.Fatalf("unexpected round trip result: %+v", out)
	}
}

func TestLoadMissingIsEmpty(t *testing.T) {
	s := newTestStore(t)
	var out []record
	if err := s.Load("nope", &out); err != nil {
		t.Fatalf("missing collection should not error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil slice, got %+v", out)
	}
}

func TestLoadEmptyFileIsEmpty(t *testing.T) {
	s := newTestStore(t)
	if err := os.WriteFile(filepath.Join(s.Dir(), "empty.json"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	var out []record
	if err := s.Load("empty", &out); err != nil {
		t.Fatalf("empty file should not error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil, got %+v", out)
	}
}

func TestSaveCreatesBackup(t *testing.T) {
	s := newTestStore(t)
	if err := s.Save("c", []record{{1, "first"}}); err != nil {
		t.Fatal(err)
	}
	if err := s.Save("c", []record{{1, "second"}}); err != nil {
		t.Fatal(err)
	}
	// 备份文件位于 <dir>/c.json.bak（不是集合 "c.bak"）
	data, err := os.ReadFile(filepath.Join(s.Dir(), "c.json.bak"))
	if err != nil {
		t.Fatalf("backup file should exist: %v", err)
	}
	var bak []record
	if err := json.Unmarshal(data, &bak); err != nil {
		t.Fatal(err)
	}
	if len(bak) != 1 || bak[0].Name != "first" {
		t.Fatalf("backup should hold previous version, got %+v", bak)
	}
}

func TestWithLockSerializesIncrements(t *testing.T) {
	s := newTestStore(t)
	const n = 100
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.WithLock("c", func() error {
				var cur []record
				if err := s.Load("c", &cur); err != nil {
					return err
				}
				cur = append(cur, record{ID: len(cur)})
				return s.Save("c", cur)
			})
		}()
	}
	wg.Wait()

	var out []record
	if err := s.Load("c", &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != n {
		t.Fatalf("lost updates: want %d records, got %d", n, len(out))
	}
}

func TestConcurrentLoadNoPartialRead(t *testing.T) {
	s := newTestStore(t)
	if err := s.Save("c", []record{{1, "a"}}); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		for i := 0; i < 500; i++ {
			_ = s.WithLock("c", func() error {
				return s.Save("c", []record{{i, "writer"}})
			})
		}
		close(done)
	}()
	for {
		var out []record
		if err := s.Load("c", &out); err != nil {
			t.Fatalf("reader saw error (partial read?): %v", err)
		}
		select {
		case <-done:
			return
		default:
		}
	}
}
