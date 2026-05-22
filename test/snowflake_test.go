package test

import (
	"ginskeleton/app/global/variable"
	_ "ginskeleton/bootstrap"
	"sync"
	"testing"
)

func TestSnowFlake(t *testing.T) {

	var slice1 []int64
	var vMuext sync.Mutex
	var wg sync.WaitGroup
	wg.Add(30000)

	for i := 1; i <= 30000; i++ {
		go func() {
			defer wg.Done()

			vMuext.Lock()
			slice1 = append(slice1, variable.SnowFlake.GetId())
			vMuext.Unlock()

		}()
	}

	wg.Wait()

	if lastLen := len(RemoveRepeatedElement(slice1)); lastLen == 30000 {
		t.Log("单元测试OK")
	} else {
		t.Errorf("雪花算法单元测试失败,并发 3万 生成的id经过去重之后，小于预期个数，去重后的个数：%d\n", lastLen)
	}
}

func RemoveRepeatedElement(arr []int64) (newArr []int64) {
	newArr = make([]int64, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
