package main

import (
	"gopkg.in/olivere/elastic.v5"
	"fmt"
	"context"
)

type LogMessage struct {
	App     string // 属于哪个项目的
	Message string // 消息内容
	Topic   string // 属于哪个topic的
}

var (
	client *elastic.Client
)

func initES(addr string) (err error) {
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("connect elasticsearch error", err)
	}

	return
}

func sendToES(topic string, data []byte) (err error) {

	msg := LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	_, err = client.Index().Index(topic).Type(topic).BodyJson(msg).Do(context.Background())
	if err != nil {
		panic(err)
		return
	}

	return
}
