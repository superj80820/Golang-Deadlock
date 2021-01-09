package main

import (
	"fmt"
	"sync"
)

type User struct {
	Balance int64
	ID      int64
	Lock    sync.RWMutex
}

func transfer(from *User, to *User, amount int64) {
	from.Lock.Lock()
	if from.Balance >= amount {
		from.Balance -= amount
		from.Lock.Unlock()
		to.Lock.Lock()
		to.Balance += amount
		to.Lock.Unlock()
	} else {
		from.Lock.Unlock()
	}
}

func main() {
	a := User{Balance: 100000000, ID: 1}
	b := User{Balance: 100000000, ID: 2}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for range [100000]int{} {
			transfer(&a, &b, 100)
		}
		wg.Done()
	}()
	go func() {
		for range [100000]int{} {
			transfer(&b, &a, 100)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
}
