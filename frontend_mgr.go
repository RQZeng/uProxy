package main

import (
	"sync"
	"errors"

	"./util"
)

var frontendMgr *FrontendMgr
var frontendMgrOnce sync.Once


func GetFrontendMgrInstance() *FrontendMgr {
	frontendMgrOnce.Do(func() {
		frontendMgr = newFrontendMgr()
	})
	return frontendMgr
}

type FrontendMgr struct {
	mMgrTbl 	map[string](*Frontend)
	mRunning	bool
}

func (this *FrontendMgr) Init(){
	this.mMgrTbl	= make(map[string](*Frontend))
	this.mRunning	= false
}

func (this *FrontendMgr) GetFrontendByID(frontendID string) (*Frontend,error) {
	f,ok := this.mMgrTbl[frontendID]
	if !ok {
		return nil,errors.New("not found frontendID")
	}
	return f,nil
}

func (this *FrontendMgr) GetFrontendByAddr(remoteAddr string) (*Frontend,error) {
	id := util.Md5Str(remoteAddr)
	return this.GetFrontendByID(id)
}

func (this *FrontendMgr) AddFrontend(f *Frontend) {
	id := f.mID
	this.mMgrTbl[id] = f
}

func (this *FrontendMgr) RemoveFrontend(f *Frontend) {
	f,ok := this.GetFrontendByID(f.mID)
	if ok != nil {
		return
	}
	delete(this.mMgrTbl,f.mID)
}

func newFrontendMgr() (*FrontendMgr) {
	m := new(FrontendMgr)
	m.Init()
	return m
}



