package app

import (
	"../conf"
	"../conn"
)

func main() {
	conf := conf.Conf{
		Address:      conf.ConnectTCPIP + ":" + conf.ConnectTCPPORT,
		MaxConnCount: 108,
		AcceptCount:  1,
	}
	server := conn.NewTcpServer(conf)
	server.Start()
}
