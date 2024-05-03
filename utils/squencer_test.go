package utils

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

var seq Sequencer

func setup() {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var nodeId int64 = 5
	seq, _ = NewSequencer(nodeId, start)
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}
func teardown() {
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetNext(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		result := map[int64]bool{}
		for i := 0; i < 100000; i++ {
			next, err := seq.Next()
			if err != nil {
				t.Fatal(err)
			}
			if _, ok := result[next.Int64()]; ok {
				t.Fatalf("Duplicate: %d", next.Int64())
			}
			result[next.Int64()] = true
		}
	})

	t.Run("multi", func(t *testing.T) {
		result := sync.Map{}
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 100000; i++ {
					next, err := seq.Next()
					if err != nil {
						t.Error(err)
					}
					if _, ok := result.Load(next.Int64()); ok {
						t.Errorf("Duplicate: %d", next.Int64())
					}
					result.Store(next.Int64(), true)
				}
			}()
		}
		wg.Wait()
	})
}

func BenchmarkSequencer(b *testing.B) {
	result := map[int64]bool{}
	for i := 0; i < b.N; i++ {
		next, err := seq.Next()
		if err != nil {
			b.Fatal(err)
		}
		if _, ok := result[next.Int64()]; ok {
			b.Fatalf("Duplicate: %d", next.Int64())
		}
		result[next.Int64()] = true
	}
}
