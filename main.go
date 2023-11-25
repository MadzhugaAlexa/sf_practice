package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func In() chan int {
	out := make(chan int)

	go func() {
		for {
			var u int
			_, err := fmt.Scanf("%d\n", &u)
			if err != nil {
				fmt.Println("Это не номер")
			} else {
				fmt.Println("Добавляем номер", u)
				out <- u
			}
		}
	}()

	return out
}

func FilterNegative(in chan int) chan int {
	out := make(chan int)

	go func() {
		for val := range in {
			if val >= 0 {
				fmt.Println("пишем положительное значение в пайп", val)
				out <- val
			}
		}
	}()

	return out
}

func FilterDivThree(in chan int) chan int {
	out := make(chan int)

	go func() {
		for val := range in {
			if val%3 != 0 {
				fmt.Println("отфильтровали значение, делящееся на 3", val)
				out <- val
			}
		}
	}()

	return out
}

func Save(b *Buffer, in chan int) {
	go func() {
		for val := range in {
			fmt.Println("сохраняем в буфер", val)
			b.Push(val)
		}
	}()
}

type Buffer struct {
	array []int
	pos   int
	size  int
	m     sync.Mutex
}

func NewBuffer(size int) *Buffer {
	return &Buffer{make([]int, size), -1, size, sync.Mutex{}}
}

func (b *Buffer) Push(val int) {
	b.m.Lock()
	defer b.m.Unlock()

	fmt.Println("добавляем в буфер", val)

	if b.pos == b.size-1 {
		for i := 1; i <= b.size-1; i++ {
			b.array[i-1] = b.array[i]
		}

		b.array[b.pos] = val
	} else {
		b.pos++
		b.array[b.pos] = val
	}
}

func (b *Buffer) Get() []int {
	b.m.Lock()
	defer b.m.Unlock()

	var output []int = b.array[:b.pos+1]
	b.pos = -1
	fmt.Println("Берем из буфера", output)
	return output
}

func Print(delay int, b *Buffer) {
	t := time.NewTicker(time.Duration(delay) * time.Second)

	go func() {
		for range t.C {
			vals := b.Get()
			if len(vals) != 0 {
				fmt.Printf("\nзначения: ")
				for _, val := range vals {
					fmt.Print(val, " ")
				}
				fmt.Println()
			}
		}
	}()
}

func Wait() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal.Aborting...\n", sig)
		os.Exit(0)
	}
}

func main() {
	delay := 10
	size := 5
	b := NewBuffer(size)

	Save(b, FilterDivThree(FilterNegative(In())))
	Print(delay, b)
	Wait()
}
