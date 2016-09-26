/*
* File: minicap.go
* Author : bigwavelet
* Description: android minicap service
* Created: 2016-09-13
 */

package minitouch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	//	"github.com/pixiv/go-libjpeg/jpeg"  // not work on windows
)

var (
	ErrAlreadyClosed = errors.New("already closed")
	HOST             = "127.0.0.1"
)

type Options struct {
	Serial string
	Port   int
	Adb    string
}

type Service struct {
	d        AdbDevice
	proc     *exec.Cmd
	port     int
	host     string
	closed   bool
	brd      *bufio.Reader
	mu       sync.Mutex
	cmdC     chan string
	restartC chan bool
}

/*
Download file to device based on goadb OpenWrite interface.
*/
func (s *Service) download(path, url string) (err error) {
	fout, err := s.d.Device.OpenWrite(path, 0755, time.Now())
	if err != nil {
		return
	}
	defer fout.Close()
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(fout, response.Body)
	if err != nil {
		return
	}
	return
}

/*
Create Minitouch Service
Description:
Serial : device serialno
Port(default: random): minitouch service port
Adb(default: adb): adb path
Eg.
	opt := Option{}
	opt.Serial = "aaa"
	service := minitouch.NewService(opt)
*/
func NewService(option Options) (s Service, err error) {
	s = Service{}
	s.d, err = newAdbDevice(option.Serial, option.Adb)
	if option.Port == 0 {
		port, err := randPort()
		if err != nil {
			return s, errors.New("port required")
		}
		s.port = port
	} else {
		s.port = option.Port
	}
	s.host = HOST
	s.closed = true
	if err != nil {
		return
	}
	s.restartC = make(chan bool, 1)
	return
}

/*
Install minitouch on device
Eg.
	service := minitouch.NewService(opt)
	err := service.Install()
P.s.
	Install function will download files, so keep network connected.
*/
func (s *Service) Install() (err error) {
	abi, err := s.d.getProp("ro.product.cpu.abi")
	if err != nil {
		return
	}
	if isExists := s.d.isFileExists("/data/local/tmp/minitouch"); isExists {
		return
	}

	url := "https://github.com/openstf/stf/raw/master/vendor/minitouch/" + abi + "/minitouch"
	err = s.download("/data/local/tmp/minitouch", url)
	return
}

/*
Check whether minitouch is supported on the device
Eg.
	service := minitouch.NewService(opt)
	supported := service.IsSupported()

For more information, see: https://github.com/openstf/minitouch
*/
func (s *Service) IsSupported() bool {
	fileExists := s.d.isFileExists("/data/local/tmp/minitouch")
	if !fileExists {
		err := s.Install()
		if err != nil {
			return false
		}
	}
	out, err := s.d.shell("/data/local/tmp/minitouch", "-h")
	if err != nil {
		return false
	}
	supported := strings.Contains(out, "-d") && strings.Contains(out, "-n") && strings.Contains(out, "-h")
	return supported
}

/*
Uninstall minicap service
Remove minicap on the device
Eg.
	service := minicap.NewService(opt)
	err := service.Uninstall()
*/
func (s *Service) Uninstall() (err error) {
	if isExists := s.d.isFileExists("/data/local/tmp/minitouch"); isExists {
		if _, err := s.d.shell("rm", "/data/local/tmp/minitouch"); err != nil {
			return err
		}
		return
	}
	return
}

/*
Start minitouch service
*/
func (s *Service) startMinitouch() (err error) {
	if !s.IsSupported() {
		return errors.New("sorry, minitouch not supported")
	}
	s.closeMinitouch()
	s.proc = s.d.buildCommand("/data/local/tmp/minitouch")
	s.proc.Stderr = os.Stderr
	stdoutReader, err := s.proc.StdoutPipe()
	if err != nil {
		return
	}
	s.brd = bufio.NewReader(stdoutReader)
	err = s.proc.Start()
	if err != nil {
		return
	}

	time.Sleep(3 * time.Second)
	if _, err = s.d.run("forward", fmt.Sprintf("tcp:%d", s.port), "localabstract:minitouch"); err != nil {
		return
	}
	s.closed = false
	s.cmdC = make(chan string, 1)
	return
}

/*
Close Minitouch stream
Eg.
	err := service.Close()
*/
func (s *Service) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrAlreadyClosed
	}
	s.closed = true
	s.closeMinitouch()
	s.d.run("forward", "--remove", fmt.Sprintf("tcp:%d", s.port))
	return
}

func (s *Service) closeMinitouch() (err error) {
	if s.proc != nil && s.proc.Process != nil {
		s.proc.Process.Signal(syscall.SIGTERM)
	}
	//kill minitouch proc on device
	err = s.d.killProc("minitouch")
	return
}

/*
check whether the minitouch stream is closed.
Eg.
	err := service.IsClosed()
*/
func (s *Service) IsClosed() bool {
	return s.closed
}

func (s *Service) sendMinitouch() (err error) {
	if err = s.startMinitouch(); err != nil {
		return
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return
	}
	go func() {
		for cmd := range s.cmdC {
			if cmd == "" {
				continue
			} else if cmd[len(cmd)-1] != '\n' {
				cmd += "\n"
			}
			_, err := conn.Write([]byte(cmd))
			if err != nil {
				log.Fatal(err)
				s.restartC <- true
			}
		}
		conn.Close()
	}()
	return
}

func (s *Service) Start() (err error) {
	s.sendMinitouch()
	return
	if err = s.sendMinitouch(); err != nil {
		return
	}

	go func() {
		for {
			<-s.restartC
			if err := s.startMinitouch(); err != nil {
				break
			}
		}
	}()
	return
}

/*
Click postiion(x,y)
*/
func (s *Service) Click(x, y int) {
	cmd := fmt.Sprintf("d 0 %d %d 50\nc\nu 0\nc\n", x, y)
	s.cmdC <- cmd
	return
}

/*
Swipe from (sx, sy) to (ex, ey)
*/

func (s *Service) Swipe(sx, sy, ex, ey int) {
	step := 10
	dx := (ex - sx) / step
	dy := (ey - sy) / step
	s.cmdC <- fmt.Sprintf("d 0 %d %d 50\nc\n", sx, sy)
	for i := 0; i < step; i++ {
		x, y := sx+i*dx, sy+i*dy
		s.cmdC <- fmt.Sprintf("m 0 %d %d 50\nc\n", x, y)
	}
	s.cmdC <- fmt.Sprintf("u 0 %d %d 50\nc\nu 0\nc\n", ex, ey)
	return

}

/*
General interface: Operation
@Parameters:
	action:
		d: down
		m: move
		u: up
	index:
		input index
	PosX:
		postion x axis
	PosY:
		postion y axis
*/
func (s *Service) Operation(action string, index, posX, posY int) {
	switch action {
	case "d":
		s.cmdC <- fmt.Sprintf("d %v %v %v 50\nc\n", index, posX, posY)
	case "m":
		s.cmdC <- fmt.Sprintf("m %v %v %v 50\nc\n", index, posX, posY)

	case "u":
		s.cmdC <- fmt.Sprintf("u %v %v %v 50\nc\nu %v\nc\n", index, posX, posY, index)

	}
}
