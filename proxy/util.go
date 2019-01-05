package proxy

import (
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/golang/glog"
)

//go:generate stringer -type=netType $GOFILE
type netType int

const (
	QUIC netType = iota
	KCP
	TCP
)

func relay(conn1, conn2 net.Conn) {
	wg := &sync.WaitGroup{}
	exitFlag := new(int32)

	go redirect(conn2, conn1, wg, exitFlag)
	redirect(conn1, conn2, wg, exitFlag)

}

func redirect(dst, src net.Conn, wg *sync.WaitGroup, exitFlag *int32) {
	if _, err := io.Copy(dst, src); err != nil && (atomic.LoadInt32(exitFlag) == 0) {
		glog.V(1).Infof("%s<>%s -> %s<>%s: %s", src.RemoteAddr(), src.LocalAddr(), dst.LocalAddr(), dst.RemoteAddr(), err)
	}
	atomic.AddInt32(exitFlag, 1)

	src.Close()
	dst.Close()
}
