package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/mqtt"
	"github.com/joek/pingu"
)

func main() {
	gbot := gobot.NewGobot()

	mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://0.0.0.0:1883", "pinguClient")
	tux := pingu.NewTux("/dev/cu.usbmodem1411")

	work := func() {
		mqttAdaptor.On("pingu", func(data []byte) {
			fmt.Println("PINGU!!!")
			tux.Wave()
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	go tux.Run("5")
	go tux.Run("3")
	gbot.AddRobot(robot)

	gbot.Start()
}
