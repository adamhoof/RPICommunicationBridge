package main

import (
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBot := TelegramBot{}
	postgreSQLHandler := PostgreSQLHandler{}

	services := ServiceContainer{
		mqtt:       &mqttHandler,
		botHandler: &telegramBot,
		db:         &postgreSQLHandler,
	}

	services.botHandler.CreateBot()

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		services.mqtt.SetupClientOptions()
		services.mqtt.CreateClient()
		services.mqtt.ConnectClient()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		services.db.Connect()
		services.db.TestConnection()
	}(&routineSyncer)

	routineSyncer.Wait()

	menuKeyboards := MenuKeyboards{}

	menuKeyboards.AllToys(&telegramBot)
	menuKeyboards.OfficeToys(&telegramBot)
	menuKeyboards.BedroomToys(&telegramBot)

	playground := Playground{}
	toyBag := ToyBag{}

	toyBag.bag = make(map[string]Toy)

	toy := postgreSQLHandler.PullToyData(1)
	toyBag.bag[toy.name] = &toy
	toyBag.bag[toy.name].Kboard(&services)

	/*var toys = []Toy{&OfficeLamp{}, &OfficeCeilLight{}, &BedroomLamp{}, &BedroomShades{}}

	toyBag.fill(toys)*/

	playground.takeOutToys(&toyBag, &services)

	telegramBot.StartBot()
}
