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
	//p.mFrontendLocalAddr = this.mLocalAddr
	//p.mFrontendRemoteAddr= this.mRemoteAddr
	l.SendTo(p)
}

func NewFrontend(remoteAddr string,localAddr string) *Frontend {
	frontend := new(Frontend)
	frontend.Init()
	frontend.mRemoteAddr	= remoteAddr
	frontend.mLocalAddr		= localAddr
	frontend.mID			= util.Md5Str(frontend.mRemoteAddr)
	return frontend
}



