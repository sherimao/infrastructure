package algorithm

import (
	"fmt"
	"sync"
)

// Function 需要并行执行的函数
type Function func(interface{}) interface{}

// Result 返回结果
type Result struct {
	Out interface{}
	Err error
}

// Pair 接收的函数及参数
type Pair struct {
	Function Function    // 函数
	In       interface{} // 入参
}

// Parallel 并行执行多个函数，按顺序返回执行结果数组
// 如果函数执行panic，则Result.Out为空，Result.Err为panic详情
// 如果函数正常执行，则Result.Out为函数执行结果，Result.Err为空
func Parallel(pairs ...Pair) (results []Result) {
	var waitGroup = sync.WaitGroup{}
	type idResult struct {
		id     int
		result Result
	}
	errorChan := make(chan idResult, len(pairs))
	safeFunction := func(i int, pair Pair) {
		defer func() {
			waitGroup.Add(-1)
			if r := recover(); r != nil {
				errorChan <- idResult{
					id:     i,
					result: Result{Err: fmt.Errorf("%v", r)},
				}
			}
		}()
		errorChan <- idResult{
			id:     i,
			result: Result{Out: pair.Function(pair.In)},
		}
	}
	for i, pair := range pairs {
		waitGroup.Add(1)
		go safeFunction(i, pair)
	}
	waitGroup.Wait()
	results = make([]Result, len(pairs))
	for len(errorChan) != 0 {
		tmp := <-errorChan
		results[tmp.id] = tmp.result
	}
	return
}
