package fsync

import (
	"fmt"
	"net"
)

type Task struct {
	Rw ReadWrite
}

func (t *Task) getServerConn(host, port string) net.Conn {
	addr, err := net.ResolveTCPAddr("tcp", ":8887")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}
	return conn

}
func (t *Task) SendFile(host, port, src, dst string) {

	conn := t.getServerConn(host, port)
	defer MyClose(conn)

	h := Handler{Rw: ReadWrite{Conn: conn}}
	h.setMsgType(MsgUploadFile)
	req := UploadReq{FName: dst}
	h.setReq(req)
	h.Rw.sendf(src)
	var resp Response
	h.getResponse(&resp)
	fmt.Println(resp.Msg)

}

func (t *Task) RecvFile(host, port, remote, local string) {

	conn := t.getServerConn(host, port)
	defer MyClose(conn)

	h := Handler{Rw: ReadWrite{Conn: conn}}
	h.setMsgType(MsgDownloadFile)
	req := DownloadReq{FName: remote}
	h.setReq(req)

	var resp Response
	h.getResponse(&resp)  //确定文件是否正常
	if resp.Code != 0{
		fmt.Println(resp.Msg)
	}else{
		h.Rw.recvf(local)
		fmt.Println("download file success")
	}
}


func (t *Task) SendDir(host, port, local, remote string) {

	itms := walk(local)
	for _, name := range itms{
		t.SendFile(host,port, name, remote + name[len(local):])
	}

}
