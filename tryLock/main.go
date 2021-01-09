package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const mutexLocked = 1 << iota

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}

type User struct {
	Balance int64
	ID      int64
	Lock    Mutex
}

func transfer(from *User, to *User, amount int64) {
	for {
		formLock := from.Lock.TryLock()
		toLock := to.Lock.TryLock()
		if formLock && toLock {
			break
		}
		if formLock {
			from.Lock.Unlock()
		}
		if toLock {
			to.Lock.Unlock()
		}
		time.Sleep(1)
	}

	if from.Balance >= amount {
		from.Balance -= amount
		to.Balance += amount
	}
	from.Lock.Unlock()
	to.Lock.Unlock()
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
