package main

import (
	//sys package
	"fmt"
	"os"
	"time"
	//"strconv"
	//"sync"
	//"net"
	//"reflect"
	"container/list"
	"flag"

	//user package
	//"./util"
	"./glog"
)

var g_quit bool 	= false
var listen_port1 *uint	= flag.Uint("listen_port1", 0, "listen port")
var listen_port2 *uint	= flag.Uint("listen_port2", 0, "listen port")
var listen_port3 *uint	= flag.Uint("listen_port3", 0, "listen port")
var listen_port4 *uint	= flag.Uint("listen_port4", 0, "listen port")
var listen_port5 *uint	= flag.Uint("listen_port5", 0, "listen port")

var svr_addr1 *string	= flag.String("svr_addr1", "192.168.0.55:50005", "svr addr")
var svr_addr2 *string	= flag.String("svr_addr2", "192.168.0.55:50006", "svr addr")
var svr_addr3 *string	= flag.String("svr_addr3", "192.168.0.55:50007", "svr addr")
var svr_addr4 *string	= flag.String("svr_addr4", "192.168.0.55:50008", "svr addr")
var svr_addr5 *string	= flag.String("svr_addr5", "192.168.0.55:50009", "svr addr")

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("    ", os.Args[0], "$port")
}

func flushLog(){
	interval := 3000 // 3s
	tick := time.NewTicker(time.Millisecond *time.Duration(interval))
	for !g_quit {
		glog.Flush()
		<-tick.C
	}
}

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Error("Start Server,Listen port1=",*listen_port1)
	glog.Error("Start Server,Listen port2=",*listen_port2)
	glog.Error("Start Server,Listen port3=",*listen_port3)
	glog.Error("Start Server,Listen port4=",*listen_port4)
	glog.Error("Start Server,Listen port5=",*listen_port5)

	glog.Error("Start Server,server addr1=",*svr_addr1)
	glog.Error("Start Server,server addr2=",*svr_addr2)
	glog.Error("Start Server,server addr3=",*svr_addr3)
	glog.Error("Start Server,server addr4=",*svr_addr4)
	glog.Error("Start Server,server addr5=",*svr_addr5)

	go flushLog()

	//go RunServer(*port,*channel_num)
	/*
	portList := [...]uint {
					*listen_port1,
					*listen_port2,
					*listen_port3,
					*listen_port4,
					*listen_port5}
	backendAddrList := [...]string {
					*svr_addr1,
					*svr_addr2,
					*svr_addr3,
					*svr_addr4,
					*svr_addr5}
	*/
	portList := list.New()
	backendAddrList := list.New()
	if *listen_port1 != 0 {
		portList.PushBack(*listen_port1)
		backendAddrList.PushBack(*svr_addr1)
	}

	if *listen_port2 != 0 {
		portList.PushBack(*listen_port2)
		backendAddrList.PushBack(*svr_addr2)
	}

	if *listen_port3 != 0 {
		portList.PushBack(*listen_port3)
		backendAddrList.PushBack(*svr_addr3)
	}

	if *listen_port4 != 0 {
		portList.PushBack(*listen_port4)
		backendAddrList.PushBack(*svr_addr4)
	}

	if *listen_port5 != 0 {
		portList.PushBack(*listen_port5)
		backendAddrList.PushBack(*svr_addr5)
	}



	//portArr := make([uint,
	listLen := portList.Len()
	portArr := make([]uint,listLen)
	backendAddrArr := make([]string,listLen)

	
	//for i:=0;i<listLen;i++ {
	i := 0
	for e:=portList.Front();e!=nil;e=e.Next() {
		portArr[i] = e.Value.(uint)
		i++
	}

	i = 0
	for e:=backendAddrList.Front();e!=nil;e=e.Next() {
		backendAddrArr[i] = e.Value.(string)
		i++
	}

	m := GetListenerMgrInstance()
	m.InitListener(portArr[0:listLen],backendAddrArr[0:listLen])
	m.Start()


	interval := 3000 // 3s
	tick := time.NewTicker(time.Millisecond *time.Duration(interval))
	for true {
		if g_quit {
			break
		}
		<-tick.C
	}
}


