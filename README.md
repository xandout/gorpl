# gorpl  [![Go Report Card](https://goreportcard.com/badge/github.com/xandout/gorpl)](https://goreportcard.com/report/github.com/xandout/gorpl)


A simple to use wrapper for [readline](https://github.com/chzyer/readline).

## Features

* Simple API
* Register callback functions

## Usage

`go get github.com/xandout/gorpl`

```go
package main
import (
    "github.com/xandout/gorpl"
    "fmt"
)

func main() {
    g, err := gorpl.New("> ", ";")

    if err != nil {
        log.Fatal(err)
    }

    g.AddAction("biggerize", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("you gave the wrong number of args")
		}
		fmt.Println(strings.ToUpper(args[0].(string)))
		return "", nil
    })
    g.Start()
}

```



### TODO

* Enable nested completion


