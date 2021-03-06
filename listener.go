package main

import (
	//"os"
	//"fmt"
	"time"
	"net"
	"strconv"


	"./glog"
	"./util"
)

const SOCKET_CHAN_LEN	= 32

type Listener struct {
	mID				string //唯一标识id
	mListenAddr		string
	mBackendAddr 	string
	mSock			*net.UDPConn
	mRunning		bool

	mSendPackageChan 	chan *NetPackage
	mRecvPackageChan 	chan *NetPackage
	mStopSignal			chan bool
	mIsLoopBack		bool

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

func (this *Listener) Init() {
	this.mID		 = ""
	this.mListenAddr = ""
	this.mBackendAddr= ""
	this.mSock		 = nil
	this.mRunning	 = false

	this.mSendPackageChan	= make(chan *NetPackage,SOCKET_CHAN_LEN)
	this.mRecvPackageChan 	= make(chan *NetPackage,SOCKET_CHAN_LEN)

	this.mRxTs		= 0
	this.mRxBps		= 0
	this.mRxPackNum	= 0
	this.mRxTotalPackNum	= 0

	this.mTxTs		= 0
	this.mTxBps		= 0
	this.mTxPackNum	= 0
	this.mTxTotalPackNum	= 0

	this.mStopSignal	= make(chan bool,2)
	this.mIsLoopBack	= false
}

func (this *Listener) GetId() string {
	return this.mID
}

func (this *Listener) GetListenAddr() string {
	return this.mListenAddr
}

func (this *Listener) GetBackendAddr() string {
	return this.mBackendAddr
}

func (this *Listener) InitSocket() {
	svrAddr := this.GetListenAddr()
	addr,ok := net.ResolveUDPAddr("udp", svrAddr)
	if ok != nil {
		glog.Error("ResolveUDPAddr err",ok)
		return
	}


	sock,err := net.ListenUDP("udp",addr)
	if err != nil {
		glog.Error("DialUDP err",err)
		return
	}
	this.mSock = sock
	return
}

func (this *Listener) grRecv() {
	glog.Error("ID=",this.mID," grRecv Start")
	defer this.mSock.Close()
	//BUFF_LEN := 1024*10
	for this.mRunning {
		this.mSock.SetDeadline(time.Now().Add(time.Duration(1) * time.Second))
		p := LentPackage()

		n, raddr, err := this.mSock.ReadFromUDP(p.data[0:])
		//glog.Error("grRecv n=" ,n)
		if err != nil {
			nerr,ok := err.(net.Error)
			if ok && nerr.Timeout() {
				continue
			}
			if !ok || !nerr.Timeout() {
				glog.Error("grRecv recv with",this.mListenAddr,",err=",err)
				this.OnQuit()
				return
			}
		}
		remoteAddr := raddr.String()
		p.mDataLen	= n
		p.OnRecv(PACKAGE_TYPE_FRONTEND,this.mListenAddr,remoteAddr)
		this.mRecvPackageChan <- p //package will be process in func grProcPackage
	}
	glog.Error("ID=",this.mID," grRecv End")

}

func (this *Listener) grProcPackage() {
	glog.Error("ID=",this.mID," grProcPackage Start")
	for this.mRunning {
		select {
		case p,ok := <-this.mRecvPackageChan:
			///glog.Error("grProcPackage")
			if ok {
				//glog.Error("grProcPackage,createts=",p.mCreateNs)
				if this.mIsLoopBack {
					this.SendTo(p)
				}else {
					GetPackageMgrInstance().OnFrontendRecv(p)
				}
			}
		case quit,ok := <- this.mStopSignal:
			if ok {
				if quit {
					break	
				}
			}else {
				glog.Error("grProcPackage err=",ok)
			}
		//default:
			//glog.Error("grProcPackage")
		}
	}
	glog.Error("ID=",this.mID," grProcPackage End")
}

func (this *Listener) OnRecv(buff []byte,remoteAddr string) {
}

func (this *Listener) SendTo(p *NetPackage) {
	this.mSendPackageChan <- p
}

func (this *Listener) grSend() {
	glog.Error("ID=",this.mID," grSend Start")
	for this.mRunning {
		select {
		case p,ok := <-this.mSendPackageChan:
			if ok {
				before_ns := util.NanoTimeStamp()
				//glog.Error("grSend,createts=",p.mCreateNs,",dataLen=",p.mDataLen)
				data := p.data[:(p.mDataLen)]
				//glog.Error(len(data))
				raddr,err := net.ResolveUDPAddr("udp", p.mFrontendRemoteAddr)
				if err != nil {
					glog.Error("grSend from " ,this.mListenAddr ," to ",p.mFrontendRemoteAddr," error")
					continue
				}
				this.mSock.WriteToUDP(data,raddr)
				this.OnSent(p)
				after_ns := util.NanoTimeStamp()
				const EXPIRE_NS = 10 * 1000 * 1000
				if after_ns-before_ns> EXPIRE_NS {
					glog.Error("too many time to send,elapsed=",(after_ns-before_ns)/1000,1000)
				}
			}else{
				glog.Error("grSend not ok")
			}
		case quit,ok := <- this.mStopSignal:
			if ok {
				if quit {
					break	
				}
			}else {
				glog.Error("grSend not ok,",ok)
			}

		//default:
		//	glog.Error("grSend")
		}
	}
	glog.Error("grSend end,listen=",this.mListenAddr)
	glog.Error("ID=",this.mID," grSend End")
}

func (this *Listener) OnSent(p *NetPackage) {
	GetPackageMgrInstance().OnFrontendSent(p)
}


func (this *Listener) Start() {
	this.mRunning	= true
	this.InitSocket() 
	go this.grRecv()
	go this.grProcPackage()
	go this.grSend()
}

func (this *Listener) StartAsLoopBack() {
	this.mIsLoopBack = true
	this.Start()
}

func (this *Listener) OnQuit(){
	glog.Error("unbelievable,listener=",this.mListenAddr)
	this.Stop()
}

func (this *Listener) Stop() {
	this.mRunning	= false
	this.mStopSignal <- true
	this.mStopSignal <- true
}

func NewListener(port uint,backendAddr string) *Listener {
	l := new(Listener)
	l.Init()
	l.mListenAddr	= "0.0.0.0:" + strconv.Itoa(int(port))
	l.mBackendAddr	= backendAddr
	l.mID			= util.Md5Str(l.GetListenAddr())
	return l
}

