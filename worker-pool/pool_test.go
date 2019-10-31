package workerpool

import (
	"context"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := NewTask(ctx, func(c context.Context) (interface{}, error) {
		if c != ctx {
			t.Fatal()
		}
		return nil, nil
	})

	// execute task
	task.Execute()

	if r := <-task.Result(); r.Err != nil || r.Result != nil {
		t.Fatal()
	}
}

func TestNewPool(t *testing.T) {
	pool := NewPool(nil, Option{ExpandableLimit: -1})
	if pool.opt.NumberWorker != numCPU || pool.opt.ExpandableLimit != 0 || pool.opt.ExpandedLifetime != time.Minute {
		t.Fatal()
	}
}

func TestPool(t *testing.T) {
	pool := NewPool(nil, Option{})
	pool.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tasks := make([]*Task, 1024)
		for i := range tasks {
			if i&1 == 0 {
				tasks[i] = pool.Execute(func(c context.Context) (interface{}, error) {
					time.Sleep(2 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecute(func(c context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond)
					return nil, nil
				})
			} else {
				tasks[i] = pool.ExecuteWithCtx(nil, func(c context.Context) (interface{}, error) {
					time.Sleep(2 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecuteWithCtx(nil, func(c context.Context) (interface{}, error) {
					time.Sleep(time.Millisecond)
					return nil, nil
				})
			}
		}
	}()
	wg.Wait()
	pool.Stop()
}

func TestPoolWithExpandable(t *testing.T) {
	pool := NewPool(nil, Option{ExpandableLimit: 2, ExpandedLifetime: 10 * time.Millisecond})
	pool.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tasks := make([]*Task, 256)
		for i := range tasks {
			if i&1 == 0 {
				tasks[i] = pool.Execute(func(c context.Context) (interface{}, error) {
					time.Sleep(5 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecute(func(c context.Context) (interface{}, error) {
					time.Sleep(5 * time.Millisecond)
					return nil, nil
				})
			} else {
				tasks[i] = pool.ExecuteWithCtx(nil, func(c context.Context) (interface{}, error) {
					time.Sleep(5 * time.Millisecond)
					return nil, nil
				})
				pool.TryExecuteWithCtx(nil, func(c context.Context) (interface{}, error) {
					time.Sleep(5 * time.Millisecond)
					return nil, nil
				})
			}
		}
	}()
	wg.Wait()
	pool.Stop()
}

func TestCorrectness(t *testing.T) {
	pool := NewPool(nil, Option{NumberWorker: runtime.NumCPU()})
	pool.Start()

	// Calculate (1^1 + 2^2 + 3^3 + ... + 1000000^1000000) modulo 1234567
	tasks := make([]*Task, 0, 1000000)
	for i := 1; i <= 1000000; i++ {
		task := moduloTask(context.Background(), uint(i), uint(i), 1234567)
		pool.Do(task)
		tasks = append(tasks, task)
	}

	// collect task results
	var s1, s2 uint
	for i := range tasks {
		if result := <-tasks[i].Result(); result.Err != nil {
			log.Fatal(result.Err)
		} else {
			s1 = uint((uint64(s1) + uint64(result.Result.(uint))) % 1234567)
		}
	}

	// sequential computation
	for i := 1; i <= 1000000; i++ {
		s2 = uint((uint64(s2) + uint64(modulo(uint(i), uint(i), 1234567))) % 1234567)
	}
	if s1 != s2 {
		log.Fatal(s1, s2)
	}

	pool.Stop()
}

func moduloTask(ctx context.Context, a, b, N uint) *Task {
	return NewTask(ctx, func(ctx context.Context) (interface{}, error) {
		return modulo(a, b, N), nil
	})
}

// calculate a^b MODULO N
func modulo(a, b uint, N uint) uint {
	switch b {
	case 0:
		return 1 % N
	case 1:
		return a % N
	default:
		if b&1 == 0 {
			t := modulo(a, b>>1, N)
			return uint(uint64(t) * uint64(t) % uint64(N))
		}

		t := modulo(a, b>>1, N)
		t = uint(uint64(t) * uint64(t) % uint64(N))
		return uint(uint64(a) * uint64(t) % uint64(N))
	}
}
