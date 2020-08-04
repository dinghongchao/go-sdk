package main

import (
	"fmt"

	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
)

const (
	topic = "subscribe"
)

func main() {
	configs, err := conf.ParseConfigFile("config.toml")
	if err != nil {
		fmt.Printf("parse configuration failed, err: %v\n", err)
		return
	}
	c, err := client.Dial(&configs[0])
	if err != nil {
		fmt.Printf("init client failed, err: %v\n", err)
		return
	}

	message := "hello, FISCO BCOS, I am unicast client!"
	err = c.PushTopicDataRandom(topic, []byte(message))
	if err != nil {
		fmt.Printf("PushTopicDataRandom failed, err: %v\n", err)
		return
	}
	fmt.Println("PushTopicDataRandom success")
}
