package util

import (
	"time"
	"strconv"
	"crypto/md5"
	"encoding/hex"
)

func TimeStamp() uint32 {
	ts:= time.Now().Unix()
	return uint32(ts)
}

func NanoTimeStamp() int64 {
	ns := time.Now().UnixNano()
	return ns
}

func Byte2StrHex(n byte) string {
	strn := strconv.FormatInt(int64(n),16)
	if len(strn) == 1 {
		strn = "0" + strn
	}
	strn = "0x" + strn
	return strn
}

func Bytes2Hex(data []byte) string {
	l := len(data)

	if l == 0 {
		return ""
	}

	LINE_SPLITE := 8
	LINE_NUM := 16
	DELTA := " "


	// b:="010001010001"
 	//base, _ := strconv.ParseInt(b, 2, 10) 
 	//hex := strconv.FormatInt(base, 16)

 	//base, _ := strconv.ParseInt(x, 16, 10) 
 	//strconv.Format.
	str := ""
	for i,v :=range data {
		if i > 0 && (i % LINE_SPLITE== 0) {
			str += " | "
		}
		if i>0 && (i%LINE_NUM==0){
			str +="\n"
		}

		str += Byte2StrHex(v) + DELTA
	}
	return str
}

func Md5Str(data string) string {
	md5val	:= md5.Sum([]byte(data))
	md5str	:= hex.EncodeToString(md5val[0:])
	return md5str
}
