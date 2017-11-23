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
	//"container/list"
	"flag"
	"runtime"
	"encoding/json"
	"io/ioutil"

	//user package
	//"./util"
	"./glog"
)

type Channel struct {                          
	ListenPort  uint     `json:"listen"`
	SvrAddr     string  `json:"forwardto"`
}     

var g_quit bool 	= false
var channel_conf *string 	= flag.String("channel_conf","./channel.json","channel config")
var core_num	*int		= flag.Int("core_num",1,"core num")

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
	runtime.GOMAXPROCS(*core_num)
	flag.Parse()
	defer glog.Flush()

	content, _ := ioutil.ReadFile(*channel_conf)                                  
	var channels []Channel
	err := json.Unmarshal(content, &channels) 
	if err != nil {
		glog.Error(err)
		return
	}


	go flushLog()


	//portArr := make([uint,
	listLen := len(channels)
	portArr := make([]uint,listLen)
	backendAddrArr := make([]string,listLen)

	for i:=0;i<len(channels);i++ {
		portArr[i] = channels[i].ListenPort
		backendAddrArr[i] = channels[i].SvrAddr
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


