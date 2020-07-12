package throol

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadPoolExecute(t *testing.T) {
	tp := NewThreadPool(2)
	tp.Init()

	var counter int64
	processFn := func(v []interface{}) error {
		atomic.AddInt64(&counter, v[0].(int64))
		return nil
	}

	task1 := task{
		processFn:      processFn,
		proccessFnArgs: []interface{}{int64(1)},
	}

	task2 := task{
		processFn:      processFn,
		proccessFnArgs: []interface{}{int64(2)},
	}

	tp.Add(&task1)
	tp.Add(&task2)

	tp.Wait()
	assert.Equal(t, int64(3), counter)
}
