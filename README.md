# kabt

In simple point form write down what you see is wrong with the code. Please explain WHY you think it is wrong and how you would rectify it.

## func main()

### for i := 1; i <= 1000000000

* Within this loop, we are calling `go func()` with each iterations of the for loop. But the `go func()` is accessing the `i` straight from the loop, meaning while the `go func()` is about to use `i` the value of `i` may be iterated and increased many times. So, here we are not assuring that the each time loop is iterating, we are triggering one `go func()` with the value of `i` at the time of that iteration.

* `Solution` just for this problem is to pass the value of `i` as a paramter of `go func()` as `go func(c int) { ... } (i)` and that way we are assuring that each `go func()` we are triggering in for loop has got the correct value of iterator `i` within it and the `go func(c int) { ... } (i)` is pushing the right value of `i` over to the `jobs` channel. This is the solution about the problem here, but in terms of program problem, this is still a problem to execute `1000000000 go func()` and keep holding them all.

* Another problem in this loop is that we are pushing the value of `i` over the `jobs` channel in each triggered `go func()` which is a blocking statement meaning that particular `go func()` will not complete until someone is listening to that channel and consumes the value pushed by the each `go func()`. Meaning at the end of `1000000000` iterations we have `1000000000 go func()` are in memory holding their threads and waiting for the consumption of the pushed value for the thread to be completed. System will be very very slow. And the program will only progress to consume after for loop with `1000000000` iterations finished.

* `Solution` can be to wrap `for i := 1; i <= 1000000000` itself within the `go func()` and iterate through `i` and push the value of `i` and push them throgu the `jobs` channel. That way we are pushing one `int` and holding the one `go func()` until that pushed `int` get consumed. And putting `for i := 1; i <= 1000000000` within the `go func()` will allow the code control flow of `main` to be progressed straight away before even starting the `for i := 1; i <= 1000000000`.

### close(jobs)

* After the `for i := 1; i <= 1000000000` we are calling `close(jobs)` which will straight away close the `jobs` channel after looping through the for loop. So there are chances for lots of triggered `go func()` by the for loop will try then send the value of `i` over the closed channel and that will cause the `panic` and the program will be terminated.

* `Solution` can be more than one but the quick one can be to use `defer close(jobs)` instead of `close(jobs)` which will ensure that the `jobs` channel will only be closed at the end of the `main()`. For clean code point of view you can declare `defer close(jobs)` straight after the declaration of `jobs` channel `jobs := make(chan int)`.

```
jobs := make(chan int)
defer close(jobs)

..
..

// close(jobs)
```