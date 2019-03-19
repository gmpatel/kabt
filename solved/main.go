package main

import "fmt"

func main() {
	jobs := make(chan int)
	defer close(jobs)

	results := make(chan int)

	go func() {
		for i := 1; i <= 1000; i++ {
			if i%2 == 0 {
				i += 99
			}
			fmt.Println("Pushing", i, "to the channel 'jobs'")
			jobs <- i
		}
	}()
	
	jobs2 := []int{}
	for w := 1; w < 1000; w++ {
		jobs2 = append(jobs2, w)
	}

	for i, w := range jobs2 {
		go worker(w, jobs, results)
		i = i + 1

	}
	close(results)
	var sum int32 = 0

	for r := range results {
		sum += int32(r)

	}

	fmt.Println(sum)
}

func worker(id int, jobs <-chan int, results chan<- int) {

	for j := range jobs {
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
