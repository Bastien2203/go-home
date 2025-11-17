package main

import (
	"gohome/internal/adapters/basic"
	"gohome/internal/adapters/homekit"
	"gohome/internal/api"
	"gohome/internal/parser"
	"gohome/internal/scanner"
	"gohome/internal/state"
)

func main() {
	store := state.NewMemoryStore()
	server := api.NewServer(store)

	bluetoothScanner := scanner.NewBluetoothScanner([]string{})
	httpScanner := scanner.NewHttpScanner([]string{"sensor1", "sensor2"}, server.Router)

	store.SaveParser(parser.NewBthomeParser(bluetoothScanner))
	store.SaveParser(parser.NewBasicParser(httpScanner))

	store.SaveAdapter(homekit.NewHomeKitAdapter())
	store.SaveAdapter(basic.NewBasicAdapter())

	server.Start()
}

//addr := "d69eee9c-8848-5c4c-3c41-d5b88b929976"
//name := "Temperature Sensor"
