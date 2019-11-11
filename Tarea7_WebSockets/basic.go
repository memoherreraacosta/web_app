package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func imprimirLetras() {
	for i := 'a'; i < 'a'+5; i++ {
		fmt.Printf("%c ", i)
	}
	fmt.Println("")
}

func f1() {
	wg.Add(2)
	go imprimirLetras()
	go imprimirLetras()
	wg.Wait()
}

func main() {
	f1()
}
