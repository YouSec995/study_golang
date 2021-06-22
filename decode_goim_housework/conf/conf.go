package conf

var (
	ConnectRPCIP = []string{
		"127.0.0.1:10086",
	}
	ConnectTCPIP   = "127.0.0.1"
	ConnectTCPPORT = "10088"
)

const (
	TypeLen       = 2                 // 消息类型字节数组长度
	LenLen        = 2                 // 消息长度字节数组长度
	HeadLen       = 2                 // 消息头字节数组长度
	ContentMaxLen = 4096              // 消息体最大长度
	BufLen        = ContentMaxLen + 8 // 缓冲buffer字节数组长度
)

// conf server配置文件
type Conf struct {
	Address      string // 地址
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收最大连接goroutine数量
}
