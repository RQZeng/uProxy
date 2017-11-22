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
	})
	return proxyerMgr
}

type ProxyerMgr struct {
	mMgrTbl 	map[string](*Proxyer)
	mRunning	bool
}

func (this *ProxyerMgr) Init(){
	this.mMgrTbl	= make(map[string](*Proxyer))
	this.mRunning	= false
}

func (this *ProxyerMgr) AddProxy(p *Proxyer) {
	id := p.mID
	this.mMgrTbl[id] = p
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
	c := NewConnector(backendRemoteAddr)
	GetConnectorMgrInstance().AddConnector(c)
	e := NewBackend(c.GetBackendAddr(),c.GetLocalAddr()) //后端实例
	GetBackendMgrInstance().AddBackend(e)

	//通道
	p := NewProxyer(f.mID,e.mID) //通道实例
	this.AddProxy(p)
	glog.Error("NewChannel,fronend=(",frontendRemoteAddr,"->",frontendLocalAddr,") <==>"," backend=(",c.GetLocalAddr(),"->",c.GetBackendAddr(),")")
}


func (this *ProxyerMgr) OnFrontendRecv(p *NetPackage) {
	remoteAddr := p.mFrontendRemoteAddr
	glog.Error("OnFrontendRecv,remoteAddr=",remoteAddr)

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

func newProxyerMgr() (*ProxyerMgr) {
	m := new(ProxyerMgr)
	m.Init()
	return m
}



