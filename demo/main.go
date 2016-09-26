package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"strings"

	minicap "github.com/BigWavelet/go-minicap"
	minitouch "github.com/BigWavelet/go-minitouch"
	"github.com/gorilla/websocket"
	"github.com/nfnt/resize"
)

var mt minitouch.Service
var mc minicap.Service
var imC <-chan image.Image
var upgrader = websocket.Upgrader{}

var ratio = 3

func StartMinicap() {
	var err error
	mc, err = minicap.NewService(minicap.Options{Serial: "EP7333W7XB"})
	if err != nil {
		log.Fatal(err)
	}
	err = mc.Install()
	if err != nil {
		log.Fatal(err)
	}
	imC, err = mc.Capture()
	if err != nil {
		log.Fatal(err)
	}
}

func StartMinitouch() {
	var err error
	mt, err = minitouch.NewService(minitouch.Options{Serial: "EP7333W7XB"})
	if err != nil {
		log.Fatal(err)
	}

	err = mt.Install()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(mt.IsSupported())

	err = mt.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Action
	/*	log.Println("start operating first...")
		mt.Operation("d", 0, 100, 500)
		mt.Operation("m", 0, 300, 500)
		mt.Operation("m", 0, 900, 500)
		mt.Operation("u", 0, 900, 500)
		log.Println("start operating 2nd...")
		mt.Operation("d", 0, 100, 500)
		mt.Operation("m", 0, 300, 500)
		mt.Operation("m", 0, 900, 500)
		mt.Operation("u", 0, 900, 500)*/
}

func hIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("t").ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func hImageWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade err:", err)
		return
	}
	log.Println("ws start..........")
	done := make(chan bool, 1)
	go func() {
		buf := new(bytes.Buffer)
		buf.Reset()
		lstImg, err := mc.LastScreenshot()
		if err == nil {
			size := lstImg.Bounds().Size()
			newIm := resize.Resize(uint(size.X)/uint(ratio), 0, lstImg, resize.Lanczos3)
			wr, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Println(err)
				return
			}

			if err := jpeg.Encode(wr, newIm, nil); err != nil {
				return
			}
			wr.Close()
		}
		for im := range imC {
			//log.Println("encode image")
			select {
			case <-done:
				return
			default:
			}
			size := im.Bounds().Size()
			newIm := resize.Resize(uint(size.X)/uint(ratio), 0, im, resize.Lanczos3)
			wr, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Println(err)
				break
			}

			if err := jpeg.Encode(wr, newIm, nil); err != nil {
				break
			}
			wr.Close()
		}
	}()
	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			done <- true
			break
		}
		param := string(p)
		fields := strings.Fields(param)
		if len(fields) >= 4 {
			action := fields[0]
			index, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			posX, err := strconv.Atoi(fields[2])
			if err != nil {
				continue
			}
			posY, err := strconv.Atoi(fields[3])
			if err != nil {
				continue
			}
			posX, posY = posX*ratio, posY*ratio
			mt.Operation(action, index, posX, posY)
		}
	}
}

func startWebServer(port int) {
	http.HandleFunc("/", hIndex)
	http.HandleFunc("/ws", hImageWs)
	log.Println("start webserver here...")
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func main() {
	log.Println("start minicap service...")
	StartMinicap()
	log.Println("start minitouch service...")
	StartMinitouch()
	//log.Println("start webserver...")
	startWebServer(5678)
}
