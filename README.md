
[![GoDoc](https://godoc.org/github.com/BigWavelet/go-minitouch?status.svg)](https://godoc.org/github.com/BigWavelet/go-minitouch)

This is a minitouch library written based on golang.



## Usage

you can fetch the library by
```shell
go get github.com/BigWavelet/go-minitouch
```

Main Interface
```go
NewService() //new minitouch service

Install() // install minitouch

Start() //start minitouch service

Click(x, y) //tap position (x, y)
```

For more information, please refer to the [demo](/demo/main.go)
