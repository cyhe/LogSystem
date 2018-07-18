package main

import (
	"logSystem/tailf"
	"github.com/astaxie/beego/logs"
	"time"
	"logSystem/kafka"
	"fmt"
)

func serverRun() (err error) {

	for {
		msg := tailf.GetSingleTail()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send message failed:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}

	return
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	fmt.Printf("read msg : %s, read topic:%s\n",msg.Msg,msg.Topic)
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
