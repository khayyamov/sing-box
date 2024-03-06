package db

import (
	"sync"
)

func GetDb() DbInterface {
	return Getter().Impl
}

type StringGetter struct {
	Impl DbInterface
}

var lockLng = &sync.Mutex{}

var singleInstance *StringGetter

func Getter() *StringGetter {
	if singleInstance == nil {
		lockLng.Lock()
		defer lockLng.Unlock()
		if singleInstance == nil {
			singleInstance = &StringGetter{
				Impl: &ImplementationDb{},
			}
		}
	}
	return singleInstance
}
