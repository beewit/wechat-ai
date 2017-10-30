package test

import (
	"testing"
	"time"
	"context"
)

var (
	ctxt   context.Context
	cancel context.CancelFunc
)

func TestEnum(t *testing.T) {
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
