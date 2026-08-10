package net

import (
	"fmt"
	"sync"
	"time"
)

func ThreadTest() {
	queue := make(chan int, 10)
	mu := sync.RWMutex{}
	c := 10

	// 生产者
	for i := range 10 {
		go func(i int) {
			time.Sleep(time.Second * 1)
			queue <- i * 10

			mu.Lock()
			c--

			if c == 0 {
				close(queue)
			}

			defer mu.Unlock()
		}(i)
	}

	// 消费者

	for v := range queue {
		fmt.Println("读取数据:", v)
	}
}
