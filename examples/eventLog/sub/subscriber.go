package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("parameters are not enough, example \n%s 127.0.0.1:20202 hello", os.Args[0])
	}
	endpoint := os.Args[1]
	privateKey, _ := hex.DecodeString("145e247e170ba3afd6ae97e88f00dbc976c2345d511b0f6713355d19d8b80b58")
	config := &conf.Config{IsHTTP: false, ChainID: 7, CAFile: "ca.crt", Key: "sdk.key", Cert: "sdk.crt",
		IsSMCrypto: false, GroupID: 1, PrivateKey: privateKey, NodeURL: endpoint}
	var c *client.Client
	var err error
	const (
		indent = "  "
	)
	for i := 0; i < 3; i++ {
		log.Printf("%d try to connect\n", i)
		c, err = client.Dial(config)
		if err != nil {
			log.Printf("init subscriber failed, err: %v, retrying\n", err)
			continue
		}
		break
	}
	if err != nil {
		log.Fatalf("init subscriber failed, err: %v\n", err)
	}
	var eventLogParams types.EventLogParams
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("newChannelMessage error: %v", err)
	}
	fmt.Println(id.String())
	eventLogParams.FilterID = strings.ReplaceAll(id.String(), "-", "")
	eventLogParams.FromBlock = "1"
	eventLogParams.ToBlock = "latest"
	eventLogParams.GroupID = "1"
	var topics = make([]string, 1)
	topics[0] = common.BytesToHash(crypto.Keccak256([]byte("TransferEvent(int256,string,string,uint256)"))).Hex()
	eventLogParams.Topics = topics
	var addresses = make([]string, 1)
	addresses[0] = "0xd2cf82e18f3d2c5cae0de87d29994be622f3fdd3"
	eventLogParams.Addresses = addresses

	timeout := 10 * time.Second
	queryTicker := time.NewTicker(timeout)
	defer queryTicker.Stop()
	done := make(chan bool)
	err = c.SubscribeEvent(eventLogParams, func(status int, logs []types.EventLog) {
		logRes, err := json.MarshalIndent(logs, "", indent)
		if err != nil {
			fmt.Printf("logs marshalIndent error: %v", err)
		}

		log.Printf("received: %s\n", logRes)
		log.Printf("received status: %d\n", status)
		queryTicker.Stop()
		queryTicker = time.NewTicker(timeout)
		done <- true
	})
	if err != nil {
		log.Printf("subscribe event failed, err: %v\n", err)
		return
	}

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	for {
		select {
		case <-done:
			log.Println("Done!")
			os.Exit(0)
		case <-queryTicker.C:
			log.Printf("can't receive message after 10s, %s\n", time.Now().String())
			os.Exit(1)
		case <-killSignal:
			log.Println("user exit")
			os.Exit(0)
		}
	}
}