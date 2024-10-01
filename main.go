package main

import (
	"fmt"
	"math"
)

func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func divide(done <-chan interface{}, intStream <-chan int, divisor int) <-chan float64 {
	addedStream := make(chan float64)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- float64(i) / float64(divisor): 
			}
		}
	}()
	return addedStream
}

func words(done <-chan interface{}, wordStream <-chan float64) <-chan string {
	wordChan := make(chan string)
	go func() {
		defer close(wordChan)
		for { 
			select {
			case <-done:
				return
			case val, ok := <-wordStream:
				if !ok {
					return
				}
				if math.Mod(val, 2) == 0 {
					wordChan <- fmt.Sprintf("%.2f банана получила обезьяна", val)
				} else {
					wordChan <- "%"
				}
			}
		}
	}()
	return wordChan
}

func printWords(in <-chan string) {
	for v := range in { 
		fmt.Println(v)
	}
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
			}
		}
	}()
	return addedStream
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	//  (x + 1) * 2 + 5
	divideResult := divide(done,intStream,2)
	wordsResult := words(done,divideResult)
	printWords(wordsResult)
	addResult1 := add(done, intStream, 1)
	multResult := multiply(done, addResult1, 2)
	addResult2 := add(done, multResult, 5)
	printer(addResult2)
}
