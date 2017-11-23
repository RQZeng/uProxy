package main

import (
	//"os"
	//"fmt"
	"time"
	"net"
	_ "strconv"


	"./glog"
	"./util"
)

//const SOCKET_CHAN_LEN	= 32

type Connector struct {
	mID				string //唯一标识id
	mLocalAddr		string
	mBackendAddr 	string
	mSock			*net.UDPConn
	mRunning		bool

	mSendPackageChan 	chan *NetPackage
	mRecvPackageChan 	chan *NetPackage
	mStopSignal			chan bool

	//send
	mRxTs			uint
	mRxBps			uint
	mRxPackNum		uint
	mRxTotalPackNum	uint
	//recv
	mTxTs			uint
	mTxBps			uint
	mTxPackNum		uint
	mTxTotalPackNum	uint
}

func (this *Connector) Init() {
	this.mID		 = ""
	this.mLocalAddr = ""
	this.mBackendAddr= ""
	this.mSock		 = nil
	this.mRunning	 = false

	this.mSendPackageChan	= make(chan *NetPackage,SOCKET_CHAN_LEN)
	this.mRecvPackageChan 	= make(chan *NetPackage,SOCKET_CHAN_LEN)
	this.mStopSignal		= make(chan bool,2)

	this.mRxTs		= 0
	this.mRxBps		= 0
	this.mRxPackNum	= 0
	this.mRxTotalPackNum	= 0

	this.mTxTs		= 0
	this.mTxBps		= 0
	this.mTxPackNum	= 0
	this.mTxTotalPackNum	= 0
}

func (this *Connector) Connect() {
	svrAddr := this.GetBackendAddr()
	addr,ok := net.ResolveUDPAddr("udp", svrAddr)
	if ok != nil {
		glog.Error("ResolveUDPAddr err",ok)
		return
	}

	sock,err := net.DialUDP("udp",nil,addr)
	if err != nil {
		glog.Error("DialUDP err ",addr,err)
		return
	}
	this.mSock = sock
	this.mLocalAddr	= sock.LocalAddr().String()
	this.mID	= util.Md5Str(this.mLocalAddr)
	
}

func (this *Connector) GetId() string {
	return this.mID
}

func (this *Connector) GetLocalAddr() string {
	return this.mLocalAddr
}

func (this *Connector) GetBackendAddr() string {
	return this.mBackendAddr
}

func (this *Connector) grRecv() {
	glog.Error("ID=",this.mID," grRecv Start")
	defer this.mSock.Close()
	//BUFF_LEN := 1024*10
	for this.mRunning {
		this.mSock.SetDeadline(time.Now().Add(time.Duration(1) * time.Second))
		p := LentPackage()

		n, raddr, err := this.mSock.ReadFromUDP(p.data[0:])
		if err != nil {
			nerr,ok := err.(net.Error)
			if ok && nerr.Timeout() {
				continue
			}
			if !ok || !nerr.Timeout() {
				glog.Error("grRecv recv with",this.mLocalAddr,",err=",err)
				this.OnQuit()
				return
			}
		}
		//glog.Error("grRecv for ",this.mLocalAddr," n=",n)

		remoteAddr := raddr.String()
		p.mDataLen	= n
		p.OnRecv(PACKAGE_TYPE_BACKEND,this.mLocalAddr,remoteAddr)
		this.mRecvPackageChan <- p //package will be process in func grProcPackage
	}
	glog.Error("ID=",this.mID," grRecv End")
}

func (this *Connector) grProcPackage() {
	glog.Error("ID=",this.mID," grProcPackage Start")
	for this.mRunning {
		select {
		case p,ok := <-this.mRecvPackageChan:
			if ok {
				//glog.Error("grProcPackage,createts=",p.mCreateNs)
				GetPackageMgrInstance().OnBackendRecv(p)
			}
		//end case
		case quit,ok := <- this.mStopSignal:
			if ok {
				if quit {
					break	
				}
			}else {
				glog.Error("grProcPackage err=",ok)
			}
		}
		//end case
	}
	glog.Error("ID=",this.mID," grProcPackage End")
}


func (this *Connector) SendTo(p *NetPackage) {
	//glog.Error("SendTo")
	this.mSendPackageChan <- p
}

func (this *Connector) grSend() {
	glog.Error("ID=",this.mID," grSend Start")
	for this.mRunning {
		select {
		case p,ok := <-this.mSendPackageChan:
			if ok {
				before_ns := util.NanoTimeStamp()
				//glog.Error("grSend,createts=",p.mCreateNs,",dataLen=",p.mDataLen)
				data := p.data[:(p.mDataLen)]
				//glog.Error(len(data))
				this.mSock.Write(data)
				this.OnSent(p)
				after_ns := util.NanoTimeStamp()
				const EXPIRE_NS = 10 * 1000 * 1000
				if after_ns-before_ns> EXPIRE_NS {
					glog.Error("too many time to send,elapsed=",(after_ns-before_ns)/1000,1000)
				}
			}
		//end case
		case quit,ok := <- this.mStopSignal:
			if ok {
				if quit {
					break	
				}
			}else {
				glog.Error("grSend err=",ok)
			}
		}
		//end case
	}
	glog.Error("ID=",this.mID," grSend End")
}

func (this *Connector) OnSent(p *NetPackage) {
	GetPackageMgrInstance().OnFrontendSent(p)
}


func (this *Connector) Start() {
	this.mRunning	= true
	go this.grRecv()
	go this.grProcPackage()
	go this.grSend()
}

func (this *Connector) OnQuit(){
	p,err := GetProxyerMgrInstance().GetProxyerByBackendAddr(this.mLocalAddr)
	if err != nil {
		glog.Error("Unbelievable not found proxyer for ",this.mLocalAddr)
		return
	}
	p.OnChannelBroken(CHANNEL_BROKEN_REASON_FRONTEND_SOCKET_ERR)
}

func (this *Connector) Stop() {
	this.mRunning	= false
	this.mStopSignal <- true
	this.mStopSignal <- true

}

func NewConnector(backendAddr string) *Connector {
	c := new(Connector)
	c.Init()
	c.mBackendAddr	= backendAddr
	c.Connect()
	c.Start()
	return c
}

