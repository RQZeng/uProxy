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
}


func NewBackend(remoteAddr string,localAddr string) *Backend {
	backend := new(Backend)
	backend.Init()
	backend.mRemoteAddr	= remoteAddr
	backend.mLocalAddr	= localAddr
	backend.mID			= util.Md5Str(backend.mLocalAddr)
	return backend
}



