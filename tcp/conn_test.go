package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			time.Sleep(3*time.Second)
			conn.Close()
		}()
		//c := NewConn(conn)
		go func(c net.Conn) {
			content := make([]byte, 1024)
			for {
				n, err2 := c.Read(content)
				if err2 != nil {
					if errors.Is(err2, io.EOF) {
						fmt.Println("EOF")
					}
					if errors.Is(err2, io.ErrClosedPipe) {
						fmt.Println("ErrClosedPipe")
					}
					if errors.Is(err2, net.ErrClosed) {
						fmt.Println("ErrClosed")
					}

					panic(err2)
				}
				fmt.Println(string(content[:n]))
			}
		}(conn)
	}
}

func TestTcpClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	content := bytes.NewBuffer([]byte{})
	msg := "你好阿世界"
	if err = binary.Write(content, binary.BigEndian, uint16(len(msg))); err != nil {
		panic(err)
	}
	if err = binary.Write(content, binary.BigEndian, []byte(msg)); err != nil {
		panic(err)
	}
	_, err = conn.Write(content.Bytes())
	if err != nil {
		panic(err)
	}

	content = bytes.NewBuffer([]byte{})
	msg = "不好呀"
	if err = binary.Write(content, binary.BigEndian, uint16(len(msg))); err != nil {
		panic(err)
	}
	if err = binary.Write(content, binary.BigEndian, []byte(msg)); err != nil {
		panic(err)
	}
	_, err = conn.Write(content.Bytes())
	if err != nil {
		panic(err)
	}
	time.Sleep(5*time.Second)
	//conn.Close()
}
