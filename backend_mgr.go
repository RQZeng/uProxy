package main

import (
	"sync"
	"errors"
)

var backendMgr *BackendMgr
var backendMgrOnce sync.Once


func GetBackendMgrInstance() *BackendMgr {
	backendMgrOnce.Do(func() {
		backendMgr = newBackendMgr()
	})
	return backendMgr
}

type BackendMgr struct {
	mMgrTbl 	map[string](*Backend)
	mRunning	bool
}

func (this *BackendMgr) Init(){
	this.mMgrTbl	= make(map[string](*Backend))
	this.mRunning	= false
}

func (this *BackendMgr) GetBackend(id string) (*Backend,error) {
	b,ok := this.mMgrTbl[id]
	if !ok {
		return nil,errors.New("not found proxyer id")
	}
	return b,nil
}

func (this *BackendMgr) AddBackend(b *Backend) {
	id := b.mID
	this.mMgrTbl[id] = b
}

func (this *BackendMgr) RemoveBackend(b *Backend) {
	id := b.mID
	b,ok := this.GetBackend(id)
	if ok != nil {
		return
	}
	delete(this.mMgrTbl,id)
}

func newBackendMgr() (*BackendMgr) {
	m := new(BackendMgr)
	m.Init()
	return m
}



