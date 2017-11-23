package main

import (
	"errors"
	"sync"
	//"container/list"

	"./glog"
	"./util"
)


var listener_mgr *ListenerMgr
var listenerMgrOnce sync.Once


func GetListenerMgrInstance() *ListenerMgr {
	listenerMgrOnce.Do(func() {
		listener_mgr = newListenerMgr()
	})
	return listener_mgr
}

type ListenerMgr struct {
	mMgrTbl map[string](*Listener)
}

func (this *ListenerMgr) Init() {
	this.mMgrTbl = make(map[string](*Listener))
}

func (this *ListenerMgr) GetListenerByAddr(listenAddr string) (*Listener,error) {
	id := util.Md5Str(listenAddr)
	return this.GetListener(id)
}

func (this *ListenerMgr) GetListener(ListenerId string) (*Listener,error) {
	h,ok := this.mMgrTbl[ListenerId]
	if !ok {
		return nil,errors.New("not found Listener id")
	}
	return h,nil
}

func (this *ListenerMgr) GetBackendSvrAddrByListenerAddr(listenAddr string) (string,error) {
	id := util.Md5Str(listenAddr)
	l,err := this.GetListener(id)
	if err != nil {
		errMsg := "has not listener to listen on "+listenAddr
		return "",errors.New(errMsg)
	}

	return l.mBackendAddr,nil
}

func (this *ListenerMgr) AddListener(Listener *Listener) {
	id := Listener.GetId()
	this.mMgrTbl[id] = Listener
}

func (this *ListenerMgr) InitListener(portList []uint,backendAddrList []string) {
	for i :=0; i < len(portList); i++ {
		port := portList[i]
		backendAddr	:= backendAddrList[i]
		l := NewListener(port,backendAddr)
		this.AddListener(l)
	}
}

func (this *ListenerMgr) Start() {
	for k,v := range(this.mMgrTbl) {
		glog.Error("Start listen id=",k,",listenAddr=",v.GetListenAddr(),",svrAddr=",v.GetBackendAddr())
		go v.Start()
	}
}

func newListenerMgr() *ListenerMgr{
	m := new(ListenerMgr)
	m.Init()
	return m
}





