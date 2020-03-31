package main

import (
	"fmt"
	"sync"
	"time"
)

func test(ch chan int, cg *sync.WaitGroup) {
	for {
		select {
		case a := <-ch:
			fmt.Println(a)

		}
		fmt.Println("default")

	}

	// cg.Done()

}

var cg sync.WaitGroup

func tickerTest() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticker begin")

		}
		fmt.Println("阻塞")
	}

}
func main() {
	// for {
	// 	a := rand.Intn(3)
	// 	fmt.Println(a)
	// }
	// ch := make(chan int)
	// go test(ch, &cg)

	// ch <- 1

	// time.Sleep(time.Second * 2)
	// ch <- 1
	// time.Sleep(time.Second * 10)

	// cg.Wait()

	ch := make(chan int)
	// cg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		for i := range ch {
			fmt.Println(i)
			fmt.Println("长度", len(ch))
		}

		// cg.Done()
	}()
	ch <- 1
	fmt.Println("加入1", len(ch))
	ch <- 1
	fmt.Println("加入1", len(ch))
	ch <- 1
	fmt.Println("加入1", len(ch))
	ch <- 1
	ch <- 1

	// cg.Wait()
	time.Sleep(time.Second * 2)
	ch <- 2
	time.Sleep(time.Second * 2)

}
