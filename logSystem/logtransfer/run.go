package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"time"
)

func run() (err error) {
	// 获取topic分区数量
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		logs.Error("Failed to get the list of partitions: ", err)
		return
	}

	fmt.Println(partitionList)

	// 遍历分区
	for partition := range partitionList {
		//sarama.OffsetNewest 每次从最新的位置消费
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("failed to start consumer for partition %d: %s\n", partition, err)

		}
		defer pc.AsyncClose()

		// 每一个分区启一个goruntine去消费
		go func(pc sarama.PartitionConsumer) {

			//kafkaClient.wg.Add(1)

			for msg := range pc.Messages() {
				// 两种方式给es
				// 1. 从协程里面取了数据,扔到另外一个channel里面,es的goroutine去取
				// 2. 直接写到es里面
				//logs.Debug("partition:%d, offset:%d, key:%s, value:%s\n",
				//msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

				err := sendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("send to es faile , err: %v", err)
				}

			}
			//kafkaClient.wg.Done()
		}(pc)
	}

	time.Sleep(time.Hour)
	//kafkaClient.wg.Wait()

	return
}
