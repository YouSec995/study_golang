package kafka

type mqAddr struct {
	addr []string
}

type Mq struct {
	Key, Value []byte			// 消息队列的键值
	Topic      string
	Partition  int32
	Offset     int64
}
