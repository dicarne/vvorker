package proxy

import (
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

func InitReverseProxy(targetEndpoint string, localEndpoint string) {
	go func() {
		ln, err := net.Listen("tcp", localEndpoint)
		if err != nil {
			logrus.Println("本地监听端口时异常", err)
		}
		for {
			clientConn, err := ln.Accept()
			if err != nil {
				logrus.Println("建立与客户端连接时异常", err)
				return
			}
			go connectPipe(clientConn, targetEndpoint)
		}
	}()
}

func connectPipe(clientConn net.Conn, serverAddr string) {
	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("建立与服务端连接时异常", err)
		_ = clientConn.Close()
		return
	}
	pipe(clientConn, serverConn)
}

func pipe(src net.Conn, dest net.Conn) {
	errChan := make(chan error, 1)
	onClose := func(_ error) {
		_ = dest.Close()
		_ = src.Close()
	}
	go func() {
		_, err := io.Copy(src, dest)
		errChan <- err
		onClose(err)
	}()
	go func() {
		_, err := io.Copy(dest, src)
		errChan <- err
		onClose(err)
	}()
	<-errChan
}
