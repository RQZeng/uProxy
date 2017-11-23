package main

import (
	"sync"
	"errors"

	"./glog"
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
	this.RemoveBackendByID(b.mID)
}

func (this *BackendMgr) RemoveBackendByID(id string) {
	b,ok := this.GetBackend(id)
	if ok != nil {
		glog.Error("RemoveBackendByID err,not found id=",id)
		return
	}
	b.OnDel()
	delete(this.mMgrTbl,id)
	this.PrintInfo()
}

func (this *BackendMgr) PrintInfo() {
	glog.Error("BackendMgr len=",len(this.mMgrTbl))
	for k,_ := range(this.mMgrTbl) {
		glog.Error("backendMgr key=",k)
	}
}

func newBackendMgr() (*BackendMgr) {
	m := new(BackendMgr)
	m.Init()
	return m
}



