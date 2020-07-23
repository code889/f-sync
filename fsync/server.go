package fsync

import (
	"errors"
	"fmt"
	"net"
	"os"
)

//var getRootDir = func() string {
//	d, _ := os.Getwd()
//	return d + "/tmp/"
//}

type Handler struct {
	Rw ReadWrite
}

func (h *Handler) setMsgType(msgType int64) {
	h.Rw.write64(msgType)
}

func (h *Handler) setReq(req interface{}) {
	h.Rw.writeReq(req)
}
func (h *Handler) getMsgType() int64 {
	return h.Rw.read64()
}
func (h *Handler) setResponse(code int, msg string) {
	res := Response{Code: code, Msg: msg}
	h.Rw.writeReq(res)
}

func (h *Handler) getResponse(res interface{}) {
	h.Rw.readReq(res)
}

func (h *Handler) getReq(req interface{}) {
	h.Rw.readReq(req)
}

func (h *Handler) Upload() {

	var req UploadReq
	h.getReq(&req)
	MyCreateDir(req.FName)
	h.Rw.recvf(req.FName)
	msg := (fmt.Sprintf("server recv file %v success", req.FName))
	h.setResponse(0, msg)

}

func (h *Handler) Download() {
	var req DownloadReq
	h.getReq(&req)
	fmt.Println("down file: ", req.FName)
	_, err := os.Stat(req.FName)
	if os.IsNotExist(err) {
		h.setResponse(-1, "file not exists")
	} else {
		h.setResponse(0, "ready to send file")
		h.Rw.sendf(req.FName)
	}
}
func (h *Handler) Setup() {

}
func (h *Handler) Process() {

	defer h.Finish()

	MsgType := h.getMsgType()
	switch MsgType {
	case MsgUploadFile:
		h.Upload()
	case MsgDownloadFile:
		h.Download()
	default:
		panic(errors.New("unknow msg type"))
	}

}
func (h *Handler) Finish() {
	_ = h.Rw.Conn.Close()
}

type Server struct {
}

func (s *Server) run() {

	lis, err := net.Listen("tcp", ":8887")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(fmt.Sprintf("accept error %v", err))
			continue
		}
		h := Handler{Rw: ReadWrite{Conn: conn}}
		h.Process()

	}
}
