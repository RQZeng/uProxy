package main

import (
	"errors"
	"sync"
	//"container/list"

	"./glog"
	"./util"
)


var connectorMgr *ConnectorMgr
var connectorMgrOnce sync.Once


func GetConnectorMgrInstance() *ConnectorMgr {
	connectorMgrOnce.Do(func() {
		connectorMgr = newConnectorMgr()
	})
	return connectorMgr
}

type ConnectorMgr struct {
	mMgrTbl map[string](*Connector)
}

func (this *ConnectorMgr) Init() {
	this.mMgrTbl = make(map[string](*Connector))
}

func (this *ConnectorMgr) GetConnectorByAddr(localAddr string) (*Connector,error) {
	id := util.Md5Str(localAddr)
	return this.GetConnector(id)
}

func (this *ConnectorMgr) GetConnector(ConnectorId string) (*Connector,error) {
	h,ok := this.mMgrTbl[ConnectorId]
	if !ok {
		return nil,errors.New("not found Connector id")
	}
	return h,nil
}

func (this *ConnectorMgr) AddConnector(c *Connector) {
	id := c.GetId()
	this.mMgrTbl[id] = c 
}

func (this *ConnectorMgr) RemoveConnectorByAddr(localAddr string){
	id := util.Md5Str(localAddr)
	this.RemoveConnector(id)
}

func (this *ConnectorMgr) RemoveConnector(id string){
	c,ok := this.mMgrTbl[id]
	if !ok {
		return
	}

	c.Stop()
	delete(this.mMgrTbl,id)
	this.PrintInfo()
}

func (this *ConnectorMgr) PrintInfo() {
	glog.Error("ConnectorMgr len=",len(this.mMgrTbl))
	for k,_ := range(this.mMgrTbl) {
		glog.Error("ConnectorMgr key=",k)
	}
}


func (this *ConnectorMgr) Start() {
	//for k,v := range(this.mMgrTbl) {
	//}
}

func newConnectorMgr() *ConnectorMgr{
	m := new(ConnectorMgr)
	m.Init()
	return m
}





