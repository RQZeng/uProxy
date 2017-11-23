package main

import (
	"strconv"
	"time"

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
	mChannelBroken		bool
	mStopSignal			chan bool
}

func (this *Proxyer) Init() {
	this.mID				= ""
	this.mFrontendID		= ""
	this.mBackendID			= ""
	this.mCreateTs			= 0
	this.mLastSendTs		= 0
	this.mLastRecvTs		= 0
	this.mChannelBroken		= false
	this.mStopSignal		= make(chan bool,1)
}

func (this *Proxyer) OnFrontendRecv(p *NetPackage) {
	if this.mChannelBroken {
		glog.Error("Channel "+p.String()+" Brokon")
		return
	}

	f,ok := GetFrontendMgrInstance().GetFrontendByID(this.mFrontendID)
	if ok != nil {
		glog.Error("not found frontend=",this.mFrontendID)
		return
	}
	f.OnRecv(p)

	b,ok := GetBackendMgrInstance().GetBackend(this.mBackendID)
	if ok != nil {
		glog.Error("not found backend=",this.mBackendID)
		return
	}
	b.SendPackage(p)
}

func (this *Proxyer) OnBackendRecv(p *NetPackage) {
	if this.mChannelBroken {
		glog.Error("Channel "+p.String()+" Brokon")
		return
	}

	b,ok := GetBackendMgrInstance().GetBackend(this.mBackendID)
	if ok != nil {
		glog.Error("not found backend=",this.mBackendID)
		return
	}
	b.OnRecv(p)

	f,ok := GetFrontendMgrInstance().GetFrontendByID(this.mFrontendID)
	if ok != nil {
		glog.Error("not found frontend=",this.mFrontendID)
		return
	}
	f.SendPackage(p)
}

func (this *Proxyer) String() string{
	str := ""
	f,fok := GetFrontendMgrInstance().GetFrontendByID(this.mFrontendID)

	if fok!= nil {
		return str
	}
	b,bok := GetBackendMgrInstance().GetBackend(this.mBackendID)
	if bok!=nil {
		return str
	}

	str += "frontend=(" + f.mRemoteAddr +"->" + f.mLocalAddr+")"
	str += "<==>"
	str += "backend=(" + b.mLocalAddr +"->" + b.mRemoteAddr+")"
	return str
}

const(
	CHANNEL_BROKEN_REASON_NONE	= iota
	CHANNEL_BROKEN_REASON_UNKNOW
	CHANNEL_BROKEN_REASON_FRONTEND_USER_CLOSE 	//前端用户主动关闭
	CHANNEL_BROKEN_REASON_FRONTEND_SOCKET_ERR	//前端socket错误
	CHANNEL_BROKEN_REASON_FRONTEND_NONE_DATA	//前端很久没有接受到数据了
	CHANNEL_BROKEN_REASON_BACKEND_SERVER_CLOSE
	CHANNEL_BROKEN_REASON_BACKEND_SOCKET_ERR
	CHANNEL_BROKEN_REASON_BACKEND_NONE_DATA	
)

func (this *Proxyer) OnChannelBroken(reason int) {
	this.mChannelBroken	= true
	this.mStopSignal <- true
	glog.Error("OnChannelBroken "+this.String() +",reason="+strconv.Itoa(reason))

	GetProxyerMgrInstance().RemoveProxy(this.mID)
}

func (this *Proxyer) OnDel() {
	glog.Error("ID=",this.mID," OnDel")

	GetFrontendMgrInstance().RemoveFrontendByID(this.mFrontendID)
	GetBackendMgrInstance().RemoveBackendByID(this.mBackendID)

}

func (this *Proxyer) grCheckChannelStatus() {
	glog.Error("id=",this.mID," grCheckChannelStatus Start")
	const CHECK_INTERVAL_MS	= 30
	tick := time.NewTicker(time.Millisecond * time.Duration(CHECK_INTERVAL_MS))
	for !this.mChannelBroken {
		select {
		case <- tick.C:
		{
			f,fok := GetFrontendMgrInstance().GetFrontendByID(this.mFrontendID)
			b,bok := GetBackendMgrInstance().GetBackend(this.mBackendID)
			if fok == nil && bok == nil {
				if f.IsBroken() {
					this.OnChannelBroken(CHANNEL_BROKEN_REASON_FRONTEND_NONE_DATA)
					break
				}
				if b.IsBroken() {
					this.OnChannelBroken(CHANNEL_BROKEN_REASON_BACKEND_NONE_DATA)
					break
				}

			}
		}//end case
			
		case quit,ok := <-this.mStopSignal:
		{
			if ok && quit {
				glog.Error("mStopSignal")
				break
			}
		}//end case
		}
	}
	glog.Error("proxyer=",this.mID," end grCheckChannelStatus")
}

func NewProxyer(frontendID string,backendID string) *Proxyer {
	proxyer := new(Proxyer)
	proxyer.Init()
	proxyer.mFrontendID	= frontendID
	proxyer.mBackendID	= backendID
	proxyer.mID			= util.Md5Str(proxyer.mFrontendID + proxyer.mBackendID)
	go proxyer.grCheckChannelStatus()
	return proxyer
}



