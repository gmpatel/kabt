package main

import "fmt"

func main() {
	jobs := make(chan int)
	defer close(jobs)

	results := make(chan int)
	defer close(results)

	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				i += 99
			}
			jobs <- i
		}
	}()

	for w := 1; w < 10; w++ {
		go worker(jobs, results)
	}

	var sum int32
	for r := range results {
		sum += int32(r)
	}

	fmt.Println(sum)
}

func worker(jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", j)
		go func() {
			switch j % 3 {
			case 0:
				j = j * 1
			case 1:
				j = j * 2
				results <- j * 2
			case 2:
				results <- j * 3
				j = j * 3
			}
		}()
	}

}
