package main

import (
	"fmt"
	"sync"
)

func abc() {
	// runtime.GOMAXPROCS(1)
	var mu sync.Mutex
	x := make(chan int)
	for i := 0; i < 5; i++ {

		go func(ch chan int) {
			var answer int
			defer mu.Unlock()

			mu.Lock()
			fmt.Scanf("%d", &answer)
			// fmt.Println("go routine called after", p, "value of answer is", answer)

			x <- answer
			// fmt.Println("pushed answer", p, "answer is ", answer)

		}(x)

	}

	fmt.Println("value of x is 1 and", <-x)
	fmt.Println("value of x is 2 and", <-x)
	fmt.Println("value of x is 3 and", <-x)
	fmt.Println("value of x is 4 and", <-x)
	fmt.Println("value of x is 5 and", <-x)

}
