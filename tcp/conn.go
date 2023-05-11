package tcp

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	MaxBufferSize = 65535
	PacketLen = 2
)

type Conn struct {
	conn net.Conn
	rb []byte
}

func NewConn(conn net.Conn) Conn {
	return Conn{
		conn: conn,
		rb:   make([]byte, 0, MaxBufferSize),
	}
}

// Read 解决 tcp 粘包、半包问题
func (c Conn) Read() ([]byte, error) {
	// 因为每次 Read 都会把完整的数据包返回，所以每次从缓冲区头读，一定是从数据包头字节开始读

	content := make([]byte, 1024)
	// 当缓冲区数据小于数据包头长度时，先读取够数据包头长度的字节，以知道数据包长度
	if len(c.rb) < PacketLen {
		for {
			_, err := c.conn.Read(content)
			if err != nil {
				return nil, err
			}
			c.rb = append(c.rb, content...)
			if len(c.rb) >= PacketLen {
				break
			}
		}
	}
	// 从缓冲区先取出数据包头的字节，得到数据包长度
	packetLen := binary.BigEndian.Uint16(c.rb[:PacketLen])
	// 数据包长度大于缓冲区可用长度
	if packetLen > MaxBufferSize - uint16(len(c.rb)) {
		return nil, errors.New("max buffer size exceeded")
	}
	c.rb = c.rb[PacketLen:]
	// 通过 for 循环获取大于等于数据包长度的数据
	for uint16(len(c.rb)) < packetLen {
		_, err := c.conn.Read(content)
		if err != nil {
			return nil, err
		}
		c.rb = append(c.rb, content...)
		// 数据长度大于缓冲区长度
		if len(c.rb) > MaxBufferSize {
			return nil, errors.New("max buffer size exceeded")
		}
	}
	// 截取完整数据包数据返回
	result := c.rb[:packetLen]
	// 截取保留下一个数据包的数据在缓冲区
	c.rb = c.rb[packetLen:]
	return result, nil
}
