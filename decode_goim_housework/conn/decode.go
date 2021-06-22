package conn

import (
	"../conf"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type Codec struct {
	Conn    net.Conn
	ReadBuf buffer // 读缓冲
}
type buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

// Package 消息包
type Package struct {
	Code    int    // 消息类型
	Content []byte // 消息体
}

// newCodec 创建一个解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:    conn,
		ReadBuf: newBuffer(conn, conf.BufLen),
	}
}

func newBuffer(reader io.Reader, len int) buffer {
	buf := make([]byte, len)
	return buffer{reader, buf, 0, 0}
}

func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader()
}

func (b *buffer) readFromReader() (int, error) {
	b.moveByte()
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, err
}

func (b *buffer) moveByte() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

func (c *Codec) Decode() (*Package, bool) {
	typeBuf, err := c.ReadBuf.seek(0, conf.TypeLen)
	if err != nil {
		return nil, false
	}
	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(conf.TypeLen, conf.HeadLen)
	if err != nil {
		return nil, false
	}
	// 读取数据内容
	valueType := int(binary.BigEndian.Uint16(typeBuf))
	valueLen := int(binary.BigEndian.Uint16(lenBuf))
	valueBuf, err := c.ReadBuf.read(conf.HeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := Package{Code: valueType, Content: valueBuf}
	return &message, true
}

// seek 返回n个字节，而不产生移位，如果没有足够字节，返回错误
func (b *buffer) seek(start, end int) ([]byte, error) {
	if b.end-b.start >= end-start {
		buf := b.buf[b.start+start : b.start+end]
		return buf, nil
	}
	ErrNotEnough := errors.New("not enough buf")
	return nil, ErrNotEnough
}

// read 舍弃offset个字段，读取n个字段,如果没有足够的字节，返回错误
func (b *buffer) read(offset, limit int) ([]byte, error) {
	if b.len() < offset+limit {
		ErrNotEnough := errors.New("not enough buf")
		return nil, ErrNotEnough
	}
	b.start += offset
	buf := b.buf[b.start : b.start+limit]
	b.start += limit
	return buf, nil
}

func (b *buffer) len() int {
	return b.end - b.start
}
