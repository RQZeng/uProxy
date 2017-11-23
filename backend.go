package main

import (
	"./glog"
	"./util"
)

type Backend struct {
	mID					string
	mRemoteAddr			string
	mLocalAddr			string

	mCreateTs		 	uint
	mLastRecvTs			uint
	mLastSendTs			uint
}

func (this *Backend) Init() {
	this.mID				= ""
	this.mRemoteAddr		= ""
	this.mLocalAddr			= ""
	this.mCreateTs			= 0
	this.mLastSendTs		= 0
	this.mLastRecvTs		= 0
}

func (this *Backend) SendPackage(p *NetPackage) {
	c,err := GetConnectorMgrInstance().GetConnectorByAddr(this.mLocalAddr)
	if err != nil {
		glog.Error("not found connector=",this.mLocalAddr)
		return
	}
	p.OnProxy(p.mPackageType,this.mLocalAddr,this.mRemoteAddr) 
	c.SendTo(p)
	this.OnSent(p)
}

func (this *Backend) OnRecv(p *NetPackage) {
	if p == nil {
		return
	}
	this.mLastRecvTs 	= uint(util.TimeStamp())
}

func (this *Backend) OnSent(p *NetPackage) {
	if p == nil {
		return
	}
	this.mLastSendTs	= uint(util.TimeStamp())
}

func (this *Backend) IsBroken() bool {
	const EXPIRE_SEC	= 10
	curr := uint(util.TimeStamp())

	if this.mLastRecvTs + EXPIRE_SEC < curr {
		return true
	}
	return false
}


func (this *Backend) OnDel() {
	GetConnectorMgrInstance().RemoveConnectorByAddr(this.mLocalAddr)
}


func NewBackend(remoteAddr string) *Backend {
	c := NewConnector(remoteAddr)
	GetConnectorMgrInstance().AddConnector(c)

	backend := new(Backend)
	backend.Init()
	backend.mRemoteAddr	= c.GetBackendAddr()
	backend.mLocalAddr	= c.GetLocalAddr()
	backend.mID			= util.Md5Str(backend.mLocalAddr)
	return backend
}



