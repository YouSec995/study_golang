package kafka

import (
	"errors"
	"fmt"
	"github.com/YouSec995/study_golang/conf"
	"github.com/Shopify/sarama"
)

// 获取消息队列的消息
func MqGetAgtInfo() (conf.Reader, conf.MyError) {
	addr := []string{"127.0.0.1:9092"}
	consumer, err := sarama.NewConsumer(addr, nil)
	var myErr conf.MyError
	var read conf.Reader
	if err != nil {

		myErr.ErrMsg = errors.New("can not connect message queue kafka")
		myErr.ErrCode = 1
		return read, myErr
	}
	partitionList, err := consumer.Partitions("GetAgtMQ")
	if err != nil {

		myErr.ErrMsg = errors.New("fail to get list of partition of queue kafka")
		myErr.ErrCode = 2
		return read, myErr
	}
	for partition := range partitionList { // 遍历所有的分区,但作为agent的消息队列，只有一个分区
		// 针对分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("GetAgtMQ", int32(partition), sarama.OffsetNewest)
		if err != nil {
			myErr.ErrMsg = errors.New("failed to start consumer for partition of queue kafka")
			myErr.ErrCode = 3
			return read, myErr
		}
		defer pc.AsyncClose()
		// 从分区消费信息
		func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				// 日志打印
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
				read.Msg = msg.Value
				read.Index = 0
			}
		}(pc)
	}
	return read, myErr
}

// 给消息队列推消息
func MqPushAgtInfo(w conf.Writer) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	msg := &sarama.ProducerMessage{}
	msg.Topic = "PushAgtMQ"
	msg.Value = sarama.StringEncoder(w.Msg)
	addr := []string{"127.0.0.1:9092"}
	client, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		// 打印错误日志信息
		return
	}
	defer client.Close()
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		// 打印错误日志信息
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
