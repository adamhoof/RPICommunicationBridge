package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBotHandler := TelegramBotHandler{}

	var routineSyncer sync.WaitGroup

	messageProcessors := make(map[string]mqtt.MessageHandler)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		messageProcessors["tableLamp"] = TableLampMessageProcessor()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions(messageProcessors)
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		telegramBotHandler.CreateBot()
		buttons := telegramBotHandler.GenerateTableLampButtons()
		telegramBotHandler.TableLampActionHandlers(&mqttHandler, buttons)
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		postgreSQLHandler := PostgreSQLHandler{}
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.Disconnect()
	}(&routineSyncer)

	routineSyncer.Wait()

	telegramBotHandler.StartBot()
}
