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
	m.Click(120, 800)
	time.Sleep(5 * time.Second)
	// //m.Close()
}

func main() {
	test()
}
