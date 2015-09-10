package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/mqtt"
)

func main() {
	gbot := gobot.NewGobot()

	mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://0.0.0.0:1883", "pingu")

	work := func() {
		gobot.Every(5*time.Second, func() {
			mqttAdaptor.Publish("pingu", []byte("wink_l"))
		})

		gobot.Every(6*time.Second, func() {
			mqttAdaptor.Publish("pingu", []byte("wink_r"))
		})
	}

	robot := gobot.NewRobot("pinguTester",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
