package conn

import (
	"../conf"
	"log"
	"net"
	"time"
)

// ConnContext 连接上下文
type ConnContext struct {
	Codec *Codec // 编解码器
}

// TCPServer TCP服务器
type TCPServer struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的goroutine数量
}

func NewTcpServer(conf conf.Conf) *TCPServer {
	return &TCPServer{
		Address:      conf.Address,
		MaxConnCount: conf.MaxConnCount,
		AcceptCount:  conf.AcceptCount,
	}
}

// Start 启动服务器
func (t *TCPServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		log.Fatal(err)
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < t.AcceptCount; i++ {
		go t.Accept(listen)
	}
	select {}
}

// Accept 接收客户端的tcp连接
func (t *TCPServer) Accept(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}
		err = conn.SetKeepAlive(true) // 设置长连接
		if err != nil {
			log.Fatal(err)
		}
		connContext := NewConnContext(conn)
		go connContext.DealConn()
	}
}

func NewConnContext(conn *net.TCPConn) *ConnContext {
	codec := NewCodec(conn)
	return &ConnContext{Codec: codec}
}

// DealConn 处理TCP连接
func (c *ConnContext) DealConn() {
	log.Println("tcp connected")
	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(600 * time.Second))
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.Codec.Read()
		if err != nil {
			log.Fatal(err)
		}
		for {
			message, ok := c.Codec.Decode()
			if ok {
				log.Println(message)
				continue
			}
			break
		}
	}
}
