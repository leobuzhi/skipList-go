# skipList-go
a library implementing skipList by Golang

## Installing
```
go get github.com/leobuzhi/skipList-go
```

## Example
```
package  main

import (
	skipList "github.com/leobuzhi/skipList-go"
	"fmt"
)

func main(){
	intMap:= skipList.NewIntMap()
	intMap.Set(1,2)
	intMap.Set(2,4)
	intMap.Set(3,6)
	intMap.Set(4,8)
	intMap.Set(5,10)

	fmt.Println(intMap.Get(1))
	fmt.Println(intMap.Get(2))
	fmt.Println(intMap.Get(3))
	fmt.Println(intMap.Get(4))
	fmt.Println(intMap.Get(5))

	fmt.Println(intMap.Len())
}
```