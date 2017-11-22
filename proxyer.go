package main

import (
	"./glog"
	"./util"
)

type Proxyer struct {
	mID					string
	mFrontendID			string
	mBackendID			string

	mCreateTs		 	uint
	mLastRecvTs			uint
	mLastSendTs			uint
}

func (this *Proxyer) Init() {
	this.mID				= ""
	this.mFrontendID		= ""
	this.mBackendID			= ""
	this.mCreateTs			= 0
	this.mLastSendTs		= 0
	this.mLastRecvTs		= 0
}

func (this *Proxyer) OnFrontendRecv(p *NetPackage) {
	b,ok := GetBackendMgrInstance().GetBackend(this.mBackendID)
	if ok != nil {
		glog.Error("not found backend=",this.mBackendID)
		return
	}
	b.SendPackage(p)
}

func (this *Proxyer) OnBackendRecv(p *NetPackage) {
	f,ok := GetFrontendMgrInstance().GetFrontendByID(this.mFrontendID)
	if ok != nil {
		return
	}
	f.SendPackage(p)
}

func NewProxyer(frontendID string,backendID string) *Proxyer {
	proxyer := new(Proxyer)
	proxyer.Init()
	proxyer.mFrontendID	= frontendID
	proxyer.mBackendID	= backendID
	proxyer.mID			= util.Md5Str(proxyer.mFrontendID + proxyer.mBackendID)
	return proxyer
}



