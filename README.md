
[![GoDoc](https://godoc.org/github.com/BigWavelet/go-minitouch?status.svg)](https://godoc.org/github.com/BigWavelet/go-minitouch)

This is a minitouch library written based on golang.



## Usage

you can fetch the library by
```shell
go get github.com/BigWavelet/go-minitouch
```

## Main Interface
```go
NewService() //new minitouch service

Install() // install minitouch

Start() //start minitouch service

Click(x, y) //tap position (x, y)

Swipe(sx, sy, ex, ey) //swipe from (sx,sy) to (ex,ey)

/*
For more general interface, refer to Operation
@Parameters:
    action: action type, including d(press down), m(mvoe), u(press up)
    index: action index, for single touch, action should be 0; for multi-touch: index should be 0, 1, 2...
    posX: action position x axis
    posY: action posstion y axis

@ Eg.:
    m.Operation("d", 0, 100, 500)
    m.Operation("m", 0, 300, 500)
    m.Operation("m", 0, 500, 500)
    m.Operation("m", 0, 650, 500)
    m.Operation("m", 0, 750, 500)
    m.Operation("m", 0, 850, 500)
    m.Operation("m", 0, 950, 500)
    m.Operation("m", 0, 1050, 500)
    m.Operation("u", 0, 1050, 500)
*/
Operation(action, index, posX, posY)

```


For more information, please refer to the [demo](/demo/main.go)
