package main

import (
	"log"
	"time"

	minitouch "github.com/BigWavelet/go-minitouch"
)

func test() {
	m, err := minitouch.NewService(minitouch.Options{Serial: "EP7333W7XB", Port: 9797})
	if err != nil {
		log.Fatal(err)
	}

	err = m.Install()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(m.IsSupported())

	err = m.Start()
	if err != nil {
		log.Fatal(err)
	}

	// //click
	// m.Click(120, 800)
	// time.Sleep(3 * time.Second)

	// Action
	log.Println("start operating...")
	m.Operation("d", 0, 100, 500)
	m.Operation("m", 0, 300, 500)
	m.Operation("m", 0, 500, 500)
	m.Operation("m", 0, 650, 500)
	m.Operation("m", 0, 750, 500)
	m.Operation("m", 0, 850, 500)
	m.Operation("m", 0, 950, 500)
	m.Operation("m", 0, 1050, 500)
	m.Operation("u", 0, 1050, 500)

	time.Sleep(3 * time.Second)
	m.Close()
}

func main() {
	test()
}
