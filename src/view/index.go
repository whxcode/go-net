// Package view 视图类型
package view

import (
	"fmt"
)

func writeChain(queue chan<- int) {
	queue <- 10
}

func readChain(queue <-chan int) {
	v, k := <-queue

	fmt.Println(v, k)
}

func wChain(queue chan int) {
	defer func() {
		fmt.Println("execute")
	}()

	close(queue)
	close(queue)
}

func Test() {
	testChan()
}
