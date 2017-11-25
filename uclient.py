from socket import *  
import time
import threading
import struct

HOST = '127.0.0.1'  
PORT = 50001
#PORT = 50005
#PORT = 40005
BUFSIZ = 128  
ADDR = (HOST, PORT)  

udpClient = socket(AF_INET, SOCK_DGRAM)  
udpClient.settimeout(2)

df	= 'QQ'
#redunt_on = False
redunt_on = True

total		= 0
loss		= 0
fluctuate	= 0
def OnLose():
	global loss
	loss += 1


def OnFluctuate():
	global fluctuate
	fluctuate += 1


def OnTotal():
	global total
	total+= 1
	if total%10000==0:
		print 'loss=%d,fluctuate=%d,total=%d,loss_ratio=%f,fluctuate=%f' %(loss,fluctuate,total,(float(loss))/total,(float(fluctuate))/total)



def recv_loop():
	print('Start Recv,redunt_on=',redunt_on)
	last_idx = 0
	last_elapsed = 0
	recved_tbl = {} #[idx]=true
	while True:
		try:
			rdata,raddr = udpClient.recvfrom(BUFSIZ)
			i,send_ts = struct.unpack(df,rdata)
			curr_ts = long(round(time.time()*1000))
			elapsed = curr_ts-send_ts
			if redunt_on and recved_tbl.has_key(i):
				#print "HASKEY:idx=%d" %(i)
				continue

			OnTotal()
			if last_idx +1 != i:
				print "LOSS:last_idx=%d,i=%d,elapsed=%d,last_elapsed=%d" %(last_idx,i,elapsed,last_elapsed)
				OnLose()
			if (float(abs(elapsed-last_elapsed)))/max(elapsed,last_elapsed) > 0.2:
				print "FLUCTUATE:i=%d,elapsed=%d,last_elapsed=%d" %(i,elapsed,last_elapsed)
				OnFluctuate()

			if redunt_on:
				recved_tbl[i] = True
			last_idx = i
			last_elapsed = elapsed

		except timeout as e:
			pass

loop = threading.Thread(target=recv_loop,args=())
loop.start()

idx = 0
while True:  
	time.sleep(0.01)
	cts = long(round(time.time()*1000))
	data = struct.pack(df,idx,cts)
	idx += 1
	udpClient.sendto(data,ADDR)  
	if redunt_on:
		udpClient.sendto(data,ADDR)  

udpClient.close()  
