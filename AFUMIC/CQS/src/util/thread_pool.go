package util

import (
	"reflect"
	"sync"
)

type Job func()

type ThreadPool struct {
	sync.WaitGroup
}

func (t *ThreadPool) Handle(job Job) {
	t.Add(1)
	go func() {
		defer t.Done()
		job()
	}()
}

func (t *ThreadPool) HandleFunc(handle interface{}, args ...interface{}) {
	t.Add(1)
	go func() {
		defer t.Done()
		values := make([]reflect.Value, 0, len(args))
		for _, arg := range args {
			values = append(values, reflect.ValueOf(arg))
		}
		funcValue := reflect.ValueOf(handle)
		funcValue.Call(values)
	}()
}
