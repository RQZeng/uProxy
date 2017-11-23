package main

import (
	"./glog"
	"./util"
)

type Frontend struct {
	mID					string
	mRemoteAddr			string
	mLocalAddr			string

	mCreateTs		 	uint
	mLastRecvTs			uint
	mLastSendTs			uint
}

func (this *Frontend) Init() {
	this.mID				= ""
	this.mRemoteAddr		= ""
	this.mLocalAddr			= ""
	this.mCreateTs			= 0
	this.mLastSendTs		= 0
	this.mLastRecvTs		= 0
}

func (this *Frontend) SendPackage(p *NetPackage) {
	l,err := GetListenerMgrInstance().GetListenerByAddr(this.mLocalAddr)
	if err != nil {
		glog.Error("not found listen=",this.mLocalAddr)
		return
	}
	p.OnProxy(p.mPackageType,this.mLocalAddr,this.mRemoteAddr)
	l.SendTo(p)
	this.OnSent(p)
}

func (this *Frontend) OnRecv(p *NetPackage) {
	if p == nil {
		return
	}
	this.mLastRecvTs 	= uint(util.TimeStamp())
}

func (this *Frontend) OnSent(p *NetPackage) {
	if p == nil {
		return
	}
	this.mLastSendTs	= uint(util.TimeStamp())
}

func (this *Frontend) IsBroken() bool {
	const EXPIRE_SEC	= 10
	curr := uint(util.TimeStamp())

	if this.mLastRecvTs + EXPIRE_SEC < curr {
		return true
	}
	return false
}

func (this *Frontend) OnDel() {
}

func NewFrontend(remoteAddr string,localAddr string) *Frontend {
	frontend := new(Frontend)
	frontend.Init()
	frontend.mRemoteAddr	= remoteAddr
	frontend.mLocalAddr		= localAddr
	frontend.mID			= util.Md5Str(frontend.mRemoteAddr)
	return frontend
}



