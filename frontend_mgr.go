package main

import (
	"sync"
	"errors"

	"./util"
	"./glog"
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
	this.RemoveFrontendByID(f.mID)
}

func (this *FrontendMgr) RemoveFrontendByID(id string) {
	f,ok := this.GetFrontendByID(id)
	if ok != nil {
		glog.Error("RemoveFrontendByID err,not found id=",id)
		return
	}
	f.OnDel()
	delete(this.mMgrTbl,id)
	this.PrintInfo()
}

func (this *FrontendMgr) PrintInfo() {
	glog.Error("FrontendMgr len=",len(this.mMgrTbl))
	for k,_ := range(this.mMgrTbl) {
		glog.Error("FrontendMgr key=",k)
	}
}

func newFrontendMgr() (*FrontendMgr) {
	m := new(FrontendMgr)
	m.Init()
	return m
}



