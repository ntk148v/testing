package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func checkerr(err error) {
	if err == io.EOF {
		fmt.Println("eof")
		return
	}

	if err == io.ErrUnexpectedEOF {
		fmt.Println("unexpected eof")
		return
	}

	netErr, ok := err.(net.Error)
	if !ok {
		return
	}

	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return
	}

	switch t := opErr.Err.(type) {
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ETIMEDOUT:
				fmt.Println("etimedout")
			}
		}
	default:
	}

	if netErr.Timeout() {
		fmt.Println("net timeout")
		return
	}

	fmt.Println("unknow")
}

func shell(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Println(command, err)
	return err, stdout.String(), stderr.String()
}

func main() {
	// Clean iptables DPORT 12345
	shell("iptables -D INPUT -p tcp --dport 12345 -j DROP")

	ln, err := net.Listen("tcp4", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		conn.Close()
	}()

	conn, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		panic(err)
	}

	shell("iptables -A INPUT -p tcp --dport 12345 -j DROP")
	defer shell("iptables -D INPUT -p tcp --dport 12345 -j DROP")

	sc, ok := conn.(syscall.Conn)
	if !ok {
		panic("not syscall conn")
	}

	rc, err := sc.SyscallConn()
	if err != nil {
		panic(err)
	}

	err = rc.Control(func(ufd uintptr) {
		fd := int(ufd)
		cnt, _ := syscall.GetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT)
		intvl, _ := syscall.GetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL)
		idle, _ := syscall.GetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE)
		fmt.Printf("default tcp keepalive idle=%ds interval=%ds count=%d\n", idle, intvl, cnt)

		syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, 2)
		syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, 1)
		syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPIDLE, 1)
		fmt.Println("set tcp keepalive idle=1s interval=1s count=2")
	})

	if err != nil {
		panic(err)
	}

	conn.SetDeadline(time.Now().Add(1 * time.Hour))

	// wait keepalive timedout
	time.Sleep((1 + 2*1 + 1) * time.Second)

	n, err := conn.Write([]byte("abcd"))
	fmt.Println("conn send", n, err)
	checkerr(err)
}
