# go-osr

An osu! replay parser for Golang inspired by [node-osr](https://github.com/vignedev/node-osr)

## Installation
```go get github.com/robloxxa/go-osr```

## Documentation
* [GoDoc](https://pkg.go.dev/github.com/robloxxa/go-osr)

## Examples

### Read replay from file

```go
package main

import (
	"fmt"
	"github.com/robloxxa/go-osr"
)

func main() {
	r, err := goosr.NewReplayFromFile("replay.osr")
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
```
or
```go
package main

import (
	"fmt"
	"github.com/robloxxa/go-osr"
)

func main() {
	r := goosr.NewReplay()
	data := []byte{} // Some test_replays in binary
	err := r.Unmarshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
```

### Write replay data to file
```go
package main

import (
	"fmt"
	"github.com/robloxxa/go-osr"
)

func main() {
	r, err := goosr.NewReplayFromFile("replay.osr")
	if err != nil {
		panic(err)
	}
	r.CountMiss = 1 // Change some test_replays from parsed replay
	err = r.WriteToFile("replay.osr")
	if err != nil {
		panic(err)
	}
}
```
or
```go
package main

import (
	"github.com/robloxxa/go-osr"
	"io/fs"
	"os"
)

func main() {
	r, err := goosr.NewReplayFromFile("replay.osr")
	if err != nil {
		panic(err)
	}
	r.CountMiss = 1 // Change some test_replays from parsed replay
	b, err := r.Marshal()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("replay.osr", b, fs.ModePerm)
	if err != nil {
		panic(err)
	}
}
```
