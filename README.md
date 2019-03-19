# kabt | the answer to the test

![Gopher](https://static1.squarespace.com/static/5c4ea9d5697a985352030ac6/5c4ea9feaa4a99273255562d/5c4fb8c9aa4a990c63cdba02/1549821309986/4kheader.png?format=1500w)

**`Question`**

In simple point form write down what you see is wrong with the code. Please explain WHY you think it is wrong and how you would rectify it.

For the question/test code sample given please refer to the following file.

```
repository:\\question\go-test-kablamo.docx
```

## func main() | Problems

### for i := 1; i <= 1000000000 | Problems

* **`Problem & Why?`** Within this loop, we are calling `go func()` with each iterations of the for loop. But the `go func()` is accessing the `i` straight from the loop, meaning while the `go func()` is about to use `i` the value of `i` may be iterated and increased many times. So, here we are not assuring that the each time loop is iterating, we are triggering one `go func()` with the value of `i` at the time of that iteration.

* **`Solution`** just for this problem is to pass the value of `i` as a paramter of `go func()` as `go func(c int) { ... } (i)` and that way we are assuring that each `go func()` we are triggering in for loop has got the correct value of iterator `i` within it and the `go func(c int) { ... } (i)` is pushing the right value of `i` over to the `jobs` channel. This is the solution about the problem here, but in terms of program problem, this is still a problem to execute `1000000000 go func()` and keep holding them all. Also you can choose to have `buffered channels` as well with higher capacity to allow thousands of `go func()` to progress without blocking them. You can choose to declare/make `buffered channel` `jobs := make(chan int, 1000000000)`.

* **`Problem & Why?`** Another problem in this loop is that we are pushing the value of `i` over the `jobs` channel in each triggered `go func()` which is a blocking statement meaning that particular `go func()` will not complete until someone is listening to that channel and consumes the value pushed by the each `go func()`. Meaning at the end of `1000000000` iterations we have `1000000000 go func()` are in memory holding their threads and waiting for the consumption of the pushed value for the thread to be completed. System will be very very slow. And the program will only progress to consume after for loop with `1000000000` iterations finished. But disadvantage having this is you are holding lot of data in memory as well. It is better to keep using data than holding and and process all, meaning, keep collecting and keep processing data will keep the memory free and system faster. 

* **`Solution`** can be to wrap `for i := 1; i <= 1000000000` itself within the `go func()` and iterate through `i` and push the value of `i` and push them throgu the `jobs` channel. That way we are pushing one `int` and holding the one `go func()` until that pushed `int` get consumed. And putting `for i := 1; i <= 1000000000` within the `go func()` will allow the code control flow of `main` to be progressed straight away before even starting the `for i := 1; i <= 1000000000`.

### close(jobs) | Problems

* **`Problem & Why?`** After the `for i := 1; i <= 1000000000` we are calling `close(jobs)` which will straight away close the `jobs` channel after looping through the for loop. So there are chances for lots of triggered `go func()` by the for loop will try then send the value of `i` over the closed channel and that will cause the `panic` and the program will be terminated.

* **`Solution`** can be more than one but the quick and neat one can be to use `defer close(jobs)` instead of `close(jobs)` which will ensure that the `jobs` channel will only be closed at the end of the `main()`. For clean code point of view you can declare `defer close(jobs)` straight after the declaration of `jobs` channel `jobs := make(chan int)`.

```
jobs := make(chan int)
defer close(jobs)
.
..
.
// close(jobs)
```

### for i, w := range jobs2 { ... } | Problems

* **`Problem & Why?`** Having `i = i + 1` statement in this loop is not doing anything. The for loop here is to iterate over the range of `jobs2` elements and returns index `i` and the element `w` and incrementing `i` and not using it doesn't make any difference.

* **`Solution`** is to remove that statement if not useful then keep it there for confusion. It is also good for `clean code`. This will cause the `compilation error` as now we are not using the variable `i` declared by the for loop. And as we are not using the index at all to avoid the `compilation error` we can use `_` instead of `i` in for loop statement. So the for loop statement should look something like `for _, w := range jobs2 { ... }` OR you can do `for i := range jobs2 { ... }` and within for loop access the array element using the `index` `i` like `jobs2[i]` instead of `w`.

### close(results) | Problems

* **`Problem & Why?`** After the `for i, w := range jobs2 { ... }` we are calling `close(results)` which will straight away close the `results` channel after looping through the for loop. And we are starting to use `results` channel after closing the `results` channel. This will cause the `panic` and the program will be terminated.

* **`Solution`** can be more than one but the quick and neat one can be to use `defer close(results)` instead of `close(results)` which will ensure that the `results` channel will only be closed at the end of the `main()`. For clean code point of view you can declare `defer close(results)` straight after the declaration of `results` channel `results := make(chan int)`.

    ```
    results := make(chan int)
    defer close(results)
    .
    ..
    .
    // close(results)
    ```

### var sum int32 = 0 | Problems

* **`Problem & Why?`** Default value assignment is not a `compilation error` or `runtime error` or any sort of errors but it's not require at all as all golang basic variable type declarations are `zero value` by default and here `int32` is value `0` by default even without assigning `0` to it.

* **`Solution`** is to declare `sum` as `var sum int32` instead of `var sum int32 = 0`. Again this is no kind of error at all but it's good for a `clean code`.

### for w := 1; w < 1000; w++ | jobs2 := []int{} | Problems

* **`Problem & Why?`** What we are doing here is that with the `for w := 1; w < 1000; w++` loop we are generating 999 job ids and storing them to `jobs2` array and then looping again over the range of `jobs2` array and triggering `worker()` as `go routine` and passing that `jobs2` element value as `id` and we are not evenin using that `id` within `worker()` method when it is passed to it. So, looks like we just need one loop to create number of `workers` we need to create and we don't even need `id` OR `value of iterator` at all to be passed to the `worker()` method.

* **`Solution`** can be something like this...

    **`Remove the following block of code`**

    ```
    jobs2 := []int{}
    for w := 1; w < 1000; w++ {
        jobs2 = append(jobs2, w)
    }
    ```

    **`And change the following block of code`**

    ```
    for i := range jobs2 {
        go worker(jobs2[i], jobs, results)
    }
    ```

    **`To`**

    ```
    for w := 1; w < 1000; w++ {
        go worker(jobs, results)
    }
    ```

    **`And also remove the {id} parameter from the func worker()`**

    **`From`**,

    ```
    func worker(id int, jobs <-chan int, results chan<- int)
    ```

    **`To`**

    ```
    func worker(jobs <-chan int, results chan<- int)
    ```
