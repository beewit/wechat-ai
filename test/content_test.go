package test

import (
	"testing"
	"time"
	"context"
	"sync"
	"fmt"
	"runtime"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils"
)

var (
	ctxt   context.Context
	cancel context.CancelFunc
)

func TestContent(t *testing.T) {
	ctxt, cancel = context.WithCancel(context.Background())
	go show()

	go func() {
		time.Sleep(time.Second * 10)
		closeShow()
	}()

	time.Sleep(time.Second * 15)
}

func show() {
	for {
		select {
		case <-ctxt.Done():
			println("done")
			return
		default:
			println("work ......")
		}

		time.Sleep(time.Second * 1)
	}
}

func closeShow() {
	println("close 。。。。")
	cancel()
}

var counter int = 0

func TestLockSyncMutex(t *testing.T) {
	lock := &sync.Mutex{}
	for i := 0; i < 100; i++ {
		go Count(lock)
	}
	for {
		lock.Lock()
		c := counter
		lock.Unlock()

		println("c =", c)

		runtime.Gosched()

		if c >= 10 {
			break
		}
	}
}

func Count(lock *sync.Mutex) {
	lock.Lock()
	counter++
	fmt.Println("counter =", counter)
	lock.Unlock()
}

type Task struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	LastTime string `json:"last_time"`
	State    bool   `json:"state"`
}

var taskList map[string]*Task

func TestStruct(t *testing.T) {
	task := taskList["task"]
	if task == nil {
		task = new(Task)
	}
	task.Content = "张三"
	println(task.Content)
	println(convert.ToObjStr(task))
	task.Content = "张三2"
	println(convert.ToObjStr(task))

}

func TestRand(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(utils.NewRandom().Number(1))
		time.Sleep(time.Millisecond)
	}
}

func TestDefer(t *testing.T) {
	defer println("defer ... ")
	i := 0
	defterP(i)
}

func defterP(i int) {
	if i>10{
		return
	}
	i++
	println(utils.NewRandom().Number(1))
	time.Sleep(time.Second)
	defterP(i)
}
