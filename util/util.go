package util

import (
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

var Id = time.Now().Unix() * 1000

func SafeGo(fn func()) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Errorf("err recovered: %+v", e)
			}
		}()
		fn()
	}()
}

func GenerateId() int64 {
	return atomic.AddInt64(&Id, 1)
}
