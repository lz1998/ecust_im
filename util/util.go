package util

import (
	log "github.com/sirupsen/logrus"
)

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
