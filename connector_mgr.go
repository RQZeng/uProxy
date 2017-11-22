package main

import (
	"errors"
	"sync"
	//"container/list"

	_ "./glog"
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

func (this *ConnectorMgr) Start() {
	//for k,v := range(this.mMgrTbl) {
	//}
}

func newConnectorMgr() *ConnectorMgr{
	m := new(ConnectorMgr)
	m.Init()
	return m
}





