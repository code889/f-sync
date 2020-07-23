package fsync

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

type ReadWrite struct {
	Conn net.Conn
}

func (r *ReadWrite) read(size int64) []byte {
	buf := make([]byte, size, size)
	length := 0
	for {
		n, err := r.Conn.Read(buf[length:])

		if n == 0 {
			panic(errors.New("conn closed"))
		}
		if n > 0 {
			length = length + n
		}
		if err != nil {

			panic(err)
		}
		if length == len(buf) {
			break
		}
	}
	return buf
}
func (r *ReadWrite) read64() int64 {

	var n int64
	b := r.read(8)
	err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &n)
	if err != nil {
		panic(err)
	}
	return n

}
func (r *ReadWrite) readInt64() int64 {

	var n int64

	err := binary.Read(r.Conn, binary.LittleEndian, &n)
	if err != nil {
		panic(err)
	}
	return n
}

func (r *ReadWrite) write(buf []byte) {

	length := 0
	for {
		n, err := r.Conn.Write(buf[length:])
		if n > 0 {
			length = length + n
		}
		if err != nil {
			panic(err)
		}
		if length == len(buf) {
			break
		}
	}
}

func (r *ReadWrite) write64(n int64) {

	b := bytes.NewBuffer([]byte{})
	err := binary.Write(b, binary.LittleEndian, n)
	if err != nil {
		panic(err)
	}
	r.write(b.Bytes())

}

func (r *ReadWrite) writeInt64(n int64) {

	err := binary.Write(r.Conn, binary.LittleEndian, n)
	if err != nil {
		panic(err)
	}

}

func (r *ReadWrite) writeReq(req interface{}) {

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(req)
	if err != nil {
		panic(err)
	}
	r.write64(int64(len(buf.Bytes())))
	r.write(buf.Bytes())

}

func (r *ReadWrite) readReq(req interface{}) {

	n := r.read64()
	buf := r.read(n)
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	err := dec.Decode(req)
	if err != nil {
		panic(err)
	}

}

func (r *ReadWrite) sendf(path string) {

	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)

	}
	defer func() {
		_ = f.Close() // 忽略错误
	}()

	st, err := f.Stat()
	if err != nil {
		panic(err)
	}
	r.write64(st.Size())
	total := st.Size()
	var length int64
	for {
		if length == st.Size() {
			break
		}
		n, err := io.CopyN(r.Conn, f, total-length)
		if err != nil && err != io.ErrShortWrite {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		length = length + n
	}
}
func (r *ReadWrite) recvf(path string) {

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		fmt.Println("file alread exists and rewrite")
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	total := r.read64()
	var length int64
	for {
		if length == total {
			break
		}
		n, err := io.CopyN(f, r.Conn, total-length)
		if err != nil && err != io.ErrShortWrite && err != io.EOF {
			panic(err)
		}
		length = length + n
	}

}
