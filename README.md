# skipList-go
a library implementing skipList by Golang

## Installing
```
go get github.com/leobuzhi/skipList-go
```

## Example
```
package main

import (
	skipList "github.com/leobuzhi/skipList-go"
	"fmt"
)

func main() {
	intMap := skipList.NewIntMap()
	intMap.Set(1, 2)
	intMap.Set(7, 14)
	intMap.Set(5, 10)
	intMap.Set(9, 18)
	intMap.Set(3, 6)

	fmt.Println(intMap.Get(1))
	fmt.Println(intMap.Get(3))
	fmt.Println(intMap.Get(5))
	fmt.Println(intMap.Get(7))
	fmt.Println(intMap.Get(9))
	//2 true
	//6 true
	//10 true
	//14 true
	//18 true
	fmt.Println("==========")

	fmt.Println(intMap.Len())
	//5
	fmt.Println("==========")

	it := intMap.Iterator()
	for it.Next() {
		fmt.Println(it.Key(), it.Value())
	}
	//1 2
	//3 6
	//5 10
	//7 14
	//9 18
	fmt.Println("==========")

	for it.Previous(){
		fmt.Println(it.Key(), it.Value())
	}
	//7 14
	//5 10
	//3 6
	//1 2
	fmt.Println("==========")

	//左闭右开
	rangeIt:=intMap.Range(1,9)
	for rangeIt.Next(){
		fmt.Println(rangeIt.Key(), rangeIt.Value())
	}
	//1 2
	//3 6
	//5 10
	//7 14
	fmt.Println("==========")

	for rangeIt.Previous(){
		fmt.Println(rangeIt.Key(), rangeIt.Value())
	}
	//5 10
	//3 6
	//1 2
	fmt.Println("==========")

	var iterator skipList.Iterator
	iterator =intMap.Seek(3)
	fmt.Println(iterator.Key(),iterator.Value())
	iterator = intMap.Seek(6)
	fmt.Println(iterator.Key(),iterator.Value())
	iterator = intMap.SeekToFirst()
	fmt.Println(iterator.Key(),iterator.Value())
	iterator = intMap.SeekToLast()
	fmt.Println(iterator.Key(),iterator.Value())
	//3 6
	//7 14
	//1 2
	//9 18
}
```