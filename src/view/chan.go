package view

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randInt(min, max int) int {
	// 用当前时间作为种子，每次运行结果不同
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func addNum(pipe chan<- int) {
	time.Sleep(time.Second * 2)
	// 生成 0- 100 的随机数

	v := randInt(0, 100)

	pipe <- v
}

func test02() {
	var dogs []int = make([]int, 10, 10)
	dogs1 := dogs[5:6]

	for v := range 10 {
		dogs = append(dogs, v)
	}

	fmt.Println(&dogs1[0])
	fmt.Println(&dogs[5])
}

func test03() {
	f1 := func() int {
		i := 0

		defer func() {
			i++
		}()

		return i
	}

	f1()
}

func test04() {
	f1 := func() (i int) {
		defer func(i int) {
			fmt.Println("0 ")
		}(i)

		defer func(i int) {
			fmt.Println("1")
		}(i)

		defer func(i int) {
			fmt.Println("2")
		}(i)

		return 10
	}

	f1()
}

func openFile() {
	i := 0

	if i == 0 {
		panic("打开文件失败")
	}
}

func test05() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	openFile()
}

func test06() {
	queue := make(chan int, 5)
	w := sync.WaitGroup{}
	l := sync.Mutex{}

	// 生产者

	c := 20
	for i := 0; i < 20; i++ {
		go func(i int) {
			defer l.Unlock()
			queue <- i * 10

			l.Lock()
			c--
			if c == 0 {
				close(queue)
			}
		}(i)
	}

	// 消费者
	for i := 0; i < 1; i++ {
		w.Add(1)
		go func() {
			defer w.Add(-1)

			for v := range queue {
				fmt.Println("获取新数据:", v)
			}
		}()
	}

	w.Wait()

	fmt.Println("所有数据处理完成")
}

func testChan() {
	test06()
	// test05()
	// fmt.Println("testChan()")
	// test03()
	// test04()
}
