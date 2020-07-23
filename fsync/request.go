package fsync


type UploadReq struct {
	FName string // 文件名
}

type DownloadReq struct {
	FName string
}

const (
	MsgUploadFile = 128
	MsgDownloadFile = 129
)

type Response struct {
	Code int
	Msg  string
}
