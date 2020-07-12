package throol

import (
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

type task struct {
	processFn      func([]interface{}) error
	proccessFnArgs []interface{}
}

type threadPool struct {
	size int
	ch   chan *task
	wg   *sync.WaitGroup

	//stats
	errCount int64
}

func NewThreadPool(size int) threadPool {
	return threadPool{
		size: size,
		ch:   make(chan *task, size),
		wg:   &sync.WaitGroup{},
	}
}

func (tp threadPool) Init() {
	for i := 0; i < tp.size; i++ {
		go tp.execute()
	}
}

func (tp threadPool) execute() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
			tp.execute()
		}
	}()

	for {
		select {
		case task := <-tp.ch:
			func() {
				defer tp.wg.Done()

				args := task.proccessFnArgs
				fn := task.processFn
				err := fn(args)
				if err != nil {
					atomic.AddInt64(&tp.errCount, 1)
				}
			}()
		}
	}
}

func (tp threadPool) Add(task *task) {
	tp.wg.Add(1)
	tp.ch <- task
}

func (tp threadPool) Wait() {
	tp.wg.Wait()
}
