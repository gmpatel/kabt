# kabt | the answer to the test

![Gopher](https://static1.squarespace.com/static/5c4ea9d5697a985352030ac6/5c4ea9feaa4a99273255562d/5c4fb8c9aa4a990c63cdba02/1549821309986/4kheader.png?format=1500w)

**`Question`**

In simple point form write down what you see is wrong with the code. Please explain WHY you think it is wrong and how you would rectify it.

For the question/test code sample given please refer to the following file.

```
repository:\\question\go-test-kablamo.docx
```
## func worker(,,) | Problems

### for j := range jobs { ... } | Problems

* **`Problem & Why?`** Within this loop, we are calling `go func()` with each iterations of the for loop. But the `go func()` is accessing the `j` straight from the loop variable, meaning while the `go func()` is about to use `j` the value of `j` may be iterated and changed many times. So, here we are not assuring that the each time loop is iterating, we are triggering one `go func()` with the each value of `j` at the time of that iteration.

* **`Solution`** for this problem is to pass the value of `j` as a paramter of `go func()` as `go func(j int) { ... } (j)` and that way we are assuring that each `go func()` we are triggering in for loop has got the correct value of iterator `j` within it and the `go func(j int) { ... } (j)` is pushing the right value of `j`.

## func main() | Problems

### for i := 1; i <= 1000000000 { ... } | Problems

* **`Problem & Why?`** Within this loop, we are calling `go func()` with each iterations of the for loop. But the `go func()` is accessing the `i` straight from the loop, meaning while the `go func()` is about to use `i` the value of `i` may be iterated and increased many times. So, here we are not assuring that the each time loop is iterating, we are triggering one `go func()` with the value of `i` at the time of that iteration.

* **`Solution`** for this problem is to pass the value of `i` as a paramter of `go func()` as `go func(c int) { ... } (i)` and that way we are assuring that each `go func()` we are triggering in for loop has got the correct value of iterator `i` within it and the `go func(c int) { ... } (i)` is pushing the right value of `i` over to the `jobs` channel.

* **`Problem & Why?`** The long running loop is holding lot of `go func()` as we are trying to push value of `i` to the non buffered `jobs` channel. So until we start consuming the `jobs` channel all the `go func()` are on hold of waiting `jobs` channel to be available to push data.

* **`Solution`** can be to declared the `buffer channel`. Meaning instead of `jobs := make(chan int)` you can declare `jobs := make(chan int, 1000000000)`, but then if the loop counter number changes you need to update your buffer value as well. I won't use the `buffer channel` as even if we are not holding all `1000000000` data, program still reserves the memory.

* **`Problem & Why?`** We are wasting time and memory here by holding all `1000000000` values in buffered channel or holding all `go func()` in memory, because the code control won't go ahead before completing the current loop.

* **`Solution`** can be to wrap this for loop in another `go func()` as that way the code control will straight away move forward to next line where we are starting processing `jobs` channel using `worker()`. So as we start consuming data from the `jobs` channel the `go func()` or the buffered channel will keep releasing data.

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

### for w := 1; w < 1000; w++ { ... } | jobs2 := []int{} | Problems

* **`Problem & Why?`** Unnecessary double looping! What we are doing here is that with the `for w := 1; w < 1000; w++` loop we are generating 999 job ids and storing them to `jobs2` array and then looping again over the range of `jobs2` array and triggering `worker()` as `go routine` and passing that `jobs2` element value as `id` and we are not evenin using that `id` within `worker()` method when it is passed to it. So, looks like we just need one loop to create number of `workers` we need to create and we don't even need `id` OR `value of iterator` at all to be passed to the `worker()` method. Also in `for i, w := range jobs2 { ... }` we are not using `i` at all, and the instruction `i = i + 1` is not required and creates confusion.

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
    for i, w := range jobs2 {
        go worker(jobs2[i], jobs, results)
        i = i + 1
    }
    ```

    **`To`**

    ```
    for w := 1; w < 1000; w++ {
        go worker(jobs, results)
    }
    ```

    **`And also remove the {id} parameter from the func worker()as we are not using the {id} in worker() at all`**

    **`From`**,

    ```
    func worker(id int, jobs <-chan int, results chan<- int)
    ```

    **`To`**

    ```
    func worker(jobs <-chan int, results chan<- int)
    ```

### var sum int32 = 0 | Problems

* **`Problem & Why?`** Default value assignment is not a `compilation error` or `runtime error` or any sort of errors but it's not require at all as all golang basic variable type declarations are `zero value` by default and here `int32` is value `0` by default even without assigning `0` to it.

* **`Solution`** is to declare `sum` as `var sum int32` instead of `var sum int32 = 0`. Again this is no kind of error at all but it's good for a `clean code`.


### for r := range results { ... } | Problems

* **`Problem & Why?`** At the very end here we are not waiting for all the triggered `go func()` above to be finised processing and we are starting to iterate over the `results` channel. Meaning we might not have most data prepared and pushed for us for final `sum` calctulation into the channel `results`.

* **`Solution`** can be many, but simplest solution is to restructure the code to `sync` all `go func()` using `sync.WaitGroup` and `waitgroup.Wait()` before proceeding with the final bit of `sum` calculation.


