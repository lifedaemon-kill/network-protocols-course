package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var password = os.Getenv("PASSWORD")
var broker = os.Getenv("BROKER")
var username = os.Getenv("USERNAME")

const (
	studentN = "1"
)

var (
	topic1 = "user_0d181660/Student" + studentN + "/Value1"
	topic2 = "user_0d181660/Student" + studentN + "/Value2"
	topic3 = "user_0d181660/Student" + studentN + "/Value3"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(onConnectHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Ошибка подключения к брокеру: ", token.Error())
	}
	defer client.Disconnect(250)

	topics := []string{topic1, topic2, topic3}
	for _, topic := range topics {
		if token := client.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
		fmt.Printf("Клиент подписался на топик: %s\n", topic)
	}

	select {}
}

func onConnectHandler(c mqtt.Client) {
	fmt.Println("Подключено к брокеру")
}

func connectionLostHandler(c mqtt.Client, err error) {
	fmt.Printf("Соединение потеряно: %v\n", err)
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Получено сообщение из топика: %s, данные: %s\n", msg.Topic(), string(msg.Payload()))
}
