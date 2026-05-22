package container

import (
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"log"
	"strings"
	"sync"
)

var sMap sync.Map

func CreateContainersFactory() *containers {
	return &containers{}
}

type containers struct {
}

func (c *containers) Set(key string, value interface{}) (res bool) {

	if _, exists := c.KeyIsExists(key); exists == false {
		sMap.Store(key, value)
		res = true
	} else {

		if variable.ZapLog == nil {
			log.Fatal(my_errors.ErrorsContainerKeyAlreadyExists + ",请解决键名重复问题,相关键：" + key)
		} else {

			variable.ZapLog.Warn(my_errors.ErrorsContainerKeyAlreadyExists + ", 相关键：" + key)
		}
	}
	return
}

func (c *containers) Delete(key string) {
	sMap.Delete(key)
}

func (c *containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

func (c *containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

func (c *containers) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPre) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
