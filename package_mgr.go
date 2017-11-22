
package main

import (
	"sync"

	"./glog"
	"./util"
)


/// package
const PACKAGE_TYPE_FRONTEND	= 0
const PACKAGE_TYPE_BACKEND	= 1

type NetPackage struct {
	mFrontendLocalAddr 	string
	mFrontendRemoteAddr string

	mBackendLocalAddr 	string
	mBackendRemoteAddr 	string
	data				[]byte
	mDataLen			int

	mPackageType		uint //0:fronted->backend,1:backend->fronted
	mCreateNs			int64 
	mRecvDataNs			int64
	mSendDataNs			int64
	mUseTime			uint
}

func (this *NetPackage) Init() {
	const BUFF_LEN = 1024 * 4
	this.mFrontendLocalAddr		= ""
	this.mFrontendRemoteAddr 	= ""
	this.mBackendLocalAddr 		= ""
	this.mBackendRemoteAddr 	= ""
	this.data					= make([]byte,BUFF_LEN)
	this.mCreateNs				= util.NanoTimeStamp()
	this.mUseTime				= 0
}

func (this *NetPackage) Reset() {
	this.mFrontendLocalAddr		= ""
	this.mFrontendRemoteAddr 	= ""
	this.mBackendLocalAddr 		= ""
	this.mBackendRemoteAddr 	= ""
}

func (this *NetPackage) OnRecv(packageType uint,localAddr string,remoteAddr string) {
	if packageType == PACKAGE_TYPE_FRONTEND {
		this.mFrontendRemoteAddr	= remoteAddr
		this.mFrontendLocalAddr		= localAddr
	}
	if packageType == PACKAGE_TYPE_BACKEND {
		this.mBackendRemoteAddr	= remoteAddr
		this.mBackendLocalAddr		= localAddr
	}
}

func (this *NetPackage) OnProxy(packageType uint,localAddr string,remoteAddr string) {
	if packageType == PACKAGE_TYPE_BACKEND {
		this.mFrontendRemoteAddr	= remoteAddr
		this.mFrontendLocalAddr		= localAddr
	}
	if packageType == PACKAGE_TYPE_FRONTEND {
		this.mBackendRemoteAddr		= remoteAddr
		this.mBackendLocalAddr		= localAddr
	}
}

func newPackage() (*NetPackage) {
	p := new(NetPackage)
	return p
}
/// package end


//package pool
var PackagePool sync.Pool
func init() {
	PackagePool = sync.Pool{
		New: func() interface{} {
			p := newPackage()
			p.Init()
			return p
		},
	}
}

//借出
func LentPackage() (*NetPackage) {
	p := PackagePool.Get().(*NetPackage)
	p.mUseTime++
	return p
}

//归还
func ReturnPackage(p *NetPackage) {
	p.Reset()
	PackagePool.Put(p)
}
//package pool end

//package mgr
var packageMgr *PackageMgr
var packageMgrOnce sync.Once


func GetPackageMgrInstance() *PackageMgr {
	packageMgrOnce.Do(func() {
		packageMgr = newPackageMgr()
	})
	return packageMgr
}

type PackageMgr struct {
	mPackageChan 	chan *NetPackage
	mRunning		bool
	//mMgrTbl map[string](*NetPackage)
}

func (this *PackageMgr) Init(){
	const MAX_PACKAGE_NUM = 1000
	this.mPackageChan = make(chan *NetPackage,MAX_PACKAGE_NUM)
	this.mRunning	= false
}

func (this *PackageMgr) OnFrontendRecv(p *NetPackage) {
	glog.Error("OnFrontendRecv")
	this.mPackageChan <- p //package will been process in func grProcPackage
}

func (this *PackageMgr) OnFrontendSent(p *NetPackage) {
	ReturnPackage(p)
}

func (this *PackageMgr) OnBackendRecv(p *NetPackage) {
	this.mPackageChan <- p //package will been process in func grProcPackage
}

func (this *PackageMgr) OnBackendSent(p *NetPackage) {
	ReturnPackage(p)
}

func (this *PackageMgr) grProcPackage() {
	for this.mRunning {
		select {
		case p,ok := <- this.mPackageChan:
			if ok {
				glog.Error("type=",p.mPackageType)
				if p.mPackageType == PACKAGE_TYPE_FRONTEND {
					GetProxyerMgrInstance().OnFrontendRecv(p)
				}

				if p.mPackageType == PACKAGE_TYPE_BACKEND {
					GetProxyerMgrInstance().OnBackendRecv(p)
				}
			}
		}
	}
}

func (this *PackageMgr) Start() {
	this.mRunning = true
	go this.grProcPackage()
}

func (this *PackageMgr) Stop() {
	this.mRunning = false
}

func newPackageMgr() (*PackageMgr) {
	m := new(PackageMgr)
	m.Init()
	go m.Start()
	return m
}

//package mgr end

