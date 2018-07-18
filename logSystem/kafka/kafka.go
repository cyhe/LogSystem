package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(addr string) (err error) {
	// 实例化配置config
	config := sarama.NewConfig()
	// 确认收到后,会回复到ack中,确保信息安全到达
	// 发给kafka放到内存里,
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机分区,kafka分区, topic是分布式的, 可以理解为一个队列,
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 实例化生产者,指定端口
	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("init kafka producer failed,err:", err)
		return
	}

	logs.Debug("init kafka succ")

	return
}

func SendToKafka(data, topic string) (err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}


	// pid:分区id, offset:分区里的偏移量
	_, _, err = client.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v  data:%v  topic:%v", err, data, topic)

		return
	}

	//logs.Debug("send succ, pid:%v, offset:%v, topic:%v", pid, offset, topic)


	//// pid:分区id, offset:分区里的偏移量
	//pid, offset, err := client.SendMessage(msg)
	//if err != nil {
	//	logs.Error("send message failed, err:%v  data:%v  topic:%v", err, data, topic)
	//
	//	return
	//}
	//
	//logs.Debug("send succ, pid:%v, offset:%v, topic:%v", pid, offset, topic)
	return
}
