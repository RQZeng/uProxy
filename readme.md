## uProxy

### 简介
	uProxy是一个udp转发代理，使用golang语言编写
### 实现原理
### 编译安装
* 使用脚本`sh build.sh`编辑，编辑后bin文件为`uProxy`

### 使用
* 配置文件`channel.json`

		[
			{"listen":50015,"forwardto":"127.0.0.1:50005"},             
			{"listen":50115,"forwardto":"127.0.0.1:50105"},             
			{"listen":50215,"forwardto":"127.0.0.1:50205"},             
			{"listen":50315,"forwardto":"127.0.0.1:50305"},             
			{"listen":50415,"forwardto":"127.0.0.1:50405"}
		]
	* `listen`:前端配置，侦听端口
	* `forwardto`:后端配置，转发地址

* 启动：`sh s.sh`