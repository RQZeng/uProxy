package main

import (
	"sync"
	"errors"

	"./glog"
	"./util"
)

var proxyerMgr *ProxyerMgr
var proxyerMgrOnce sync.Once


func GetProxyerMgrInstance() *ProxyerMgr {
	proxyerMgrOnce.Do(func() {
		proxyerMgr = newProxyerMgr()
		proxyerMgr.Start()
	})
	return proxyerMgr
}

type ProxyerMgr struct {
	mMgrTbl 	map[string](*Proxyer)
	mRunning	bool
	mBrokenChannel chan string
}

func (this *ProxyerMgr) Init(){
	this.mMgrTbl		= make(map[string](*Proxyer))
	this.mBrokenChannel	= make(chan string,100)
	this.mRunning		= false
}

func (this *ProxyerMgr) AddProxy(p *Proxyer) {
	id := p.mID
	this.mMgrTbl[id] = p
}

func (this *ProxyerMgr) RemoveProxy(id string) {
	p,ok := this.mMgrTbl[id]
	if !ok {
		return //errors.New("not found proxyer id="+id)
	}
	p.OnDel()
	delete(this.mMgrTbl,id)

	this.PrintInfo()
}

func (this *ProxyerMgr) PrintInfo() {
	glog.Error("ProxyerMgr len=",len(this.mMgrTbl))
	for k,_ := range(this.mMgrTbl) {
		glog.Error("ProxyerMgr key=",k)
	}
}

func (this *ProxyerMgr) GetProxyerByID(proxyerID string) (*Proxyer,error) {
	p,ok := this.mMgrTbl[proxyerID]
	if !ok {
		return nil,errors.New("not found proxyer id")
	}
	return p,nil
}

func (this *ProxyerMgr) GetProxyerByFrontendAddr(remoteAddr string) (*Proxyer,error){
	id := util.Md5Str(remoteAddr)
	for _,v := range(this.mMgrTbl) {
		if v.mFrontendID == id {
			return v,nil
		}
	}
	return nil,errors.New("not found frontend id")
}

func (this *ProxyerMgr) GetProxyerByFrontendID(frontendID string) (*Proxyer,error){
	for _,v := range(this.mMgrTbl) {
		if v.mFrontendID == frontendID {
			return v,nil
		}
	}
	return nil,errors.New("not found frontend id")
}

func (this *ProxyerMgr) GetProxyerByBackendAddr(localAddr string) (*Proxyer,error){
	id := util.Md5Str(localAddr)
	for _,v := range(this.mMgrTbl) {
		if v.mBackendID == id {
			return v,nil
		}
	}
	return nil,errors.New("not found backend id")
}


func (this *ProxyerMgr) GetProxyerByBackendID(backendID string) (*Proxyer,error){
	for _,v := range(this.mMgrTbl) {
		if v.mBackendID == backendID {
			return v,nil
		}
	}
	return nil,errors.New("not found backend id")
}

func (this *ProxyerMgr) NewChannel(frontendRemoteAddr,frontendLocalAddr,backendRemoteAddr string) {

	//前端连接点
	f := NewFrontend(frontendRemoteAddr,frontendLocalAddr) //前端实例
	GetFrontendMgrInstance().AddFrontend(f)

	//后端连接点
	e := NewBackend(backendRemoteAddr) //后端实例
	GetBackendMgrInstance().AddBackend(e)

	//通道
	p := NewProxyer(f.mID,e.mID) //通道实例
	this.AddProxy(p)
	glog.Error("NewChannel,fronend=(",frontendRemoteAddr,"->",frontendLocalAddr,") <==>"," backend=(",e.mLocalAddr,"->",e.mRemoteAddr,")")
}


func (this *ProxyerMgr) OnFrontendRecv(p *NetPackage) {
	remoteAddr := p.mFrontendRemoteAddr
	//glog.Error("OnFrontendRecv,remoteAddr=",remoteAddr)

	_,err := GetFrontendMgrInstance().GetFrontendByAddr(remoteAddr)
	if err != nil {
		//创建新的通道
		localAddr := p.mFrontendLocalAddr
		svrAddr,addrErr := GetListenerMgrInstance().GetBackendSvrAddrByListenerAddr(localAddr)
		//后端节点地址
		if addrErr != nil {
			glog.Error("OnFrontendRecv,err=",addrErr)
			return
		}

		this.NewChannel(remoteAddr,localAddr,svrAddr)
	}
	proxyer,perr	:=  this.GetProxyerByFrontendAddr(remoteAddr)
	if perr != nil {
		glog.Error("OnFrontendRecv ,err=",perr)
		return
	}
	proxyer.OnFrontendRecv(p)
}

func (this *ProxyerMgr) OnBackendRecv(p *NetPackage) {
	localAddr := p.mBackendLocalAddr

	proxyer,perr	:= this.GetProxyerByBackendAddr(localAddr)
	if perr != nil {
		glog.Error("OnBackendRecv ,err=",perr)
		return
	}
	proxyer.OnBackendRecv(p)
}

func (this *ProxyerMgr) Start() {
	this.mRunning = true
}

func newProxyerMgr() (*ProxyerMgr) {
	m := new(ProxyerMgr)
	m.Init()
	return m
}



