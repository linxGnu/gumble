package queue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func testProducer(t *testing.T, queue Queue) {
	// try offer nil
	queue.Offer(nil)

	var wg sync.WaitGroup
	for i := 0; i < maxNumberProducer; i++ {
		wg.Add(1)
		go func(producer int) {
			for j := 0; j < numberEle; j++ {
				queue.Offer(&ele{key: producer, value: j})
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if queue.Peek() == nil || int(queue.Size()) != maxNumberProducer*numberEle || queue.IsEmpty() {
		t.Fatal()
	}

	m := make(map[int]map[int]struct{})
	for i := 0; i < maxNumberProducer*numberEle; i++ {
		if polled := queue.Poll(); polled == nil {
			t.Fatal()
		} else {
			e := polled.(*ele)
			if _, ok := m[e.key]; !ok {
				m[e.key] = make(map[int]struct{})
			}
			m[e.key][e.value] = struct{}{}
		}
	}

	for i := 0; i < maxNumberProducer; i++ {
		if len(m[i]) != numberEle {
			t.Fatal()
		}

		for k := range m[i] {
			if k < 0 || k >= numberEle {
				t.Fatal()
			}
		}
	}
}

func testMix(t *testing.T, q Queue, numberProducer, numberConsumer int) {
	if numberProducer > maxNumberProducer {
		numberProducer = maxNumberProducer
	}
	if numberConsumer > maxNumberConsumer {
		numberConsumer = maxNumberConsumer
	}

	for i := 0; i < numberProducer; i++ {
		go func(producer int) {
			for j := 0; j < numberEle; j++ {
				q.Offer(&ele{key: producer, value: j})
			}
		}(i)
	}

	ch := make(chan *ele, 100000)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < numberConsumer; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					if item := q.Peek(); item != nil {
						if q.IsEmpty() {
							continue
						}

						if item = q.Poll(); item != nil {
							ch <- item.(*ele)
						}
					}
				}
			}
		}(ctx)
	}

	for i := 0; i < numberConsumer; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			ct := 0
			for {
				select {
				case <-ctx.Done():
					if ct > 950 {
						fmt.Println(ct)
					}
					wg.Done()
					return
				default:
					var item interface{}
					iter := q.Iterator()

					if iter == nil {
						time.Sleep(time.Second)
						continue
					}

					counter := 0
					for iter := q.Iterator(); iter.HasNext(); {
						if item = iter.Next(); item != nil {
							if counter++; counter > 100000 {
								ct++
								break
							}
						}
					}
				}
			}
		}(ctx)
	}

	m := make(map[int]map[int]struct{})
	for i := 0; i < maxNumberProducer*numberEle; i++ {
		if polled := <-ch; polled == nil {
			t.FailNow()
		} else {
			e := (*ele)(polled)
			if _, ok := m[e.key]; !ok {
				m[e.key] = make(map[int]struct{})
			}
			m[e.key][e.value] = struct{}{}
		}
	}

	for i := 0; i < maxNumberProducer; i++ {
		if len(m[i]) != numberEle {
			t.FailNow()
		}

		for k := range m[i] {
			if k < 0 || k >= numberEle {
				t.FailNow()
			}
		}
	}
	cancel()
	wg.Wait()
}
