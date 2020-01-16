package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//ffff
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
	for {
		a := rand.Intn(3)
		fmt.Println(a)

	}
	// ch := make(chan int)
	// go test(ch, &cg)

	// ch <- 1

	// time.Sleep(time.Second * 2)
	// ch <- 1
	// time.Sleep(time.Second * 10)

	// cg.Wait()
}
