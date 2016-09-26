/*
* File: rotation.go
* Author : bigwavelet
* Description: android rotation watcher
* Created: 2016-09-13
 */

package minitouch

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Rotation struct {
	d           AdbDevice
	orientation int
	proc        *exec.Cmd
	closed      bool
	brd         *bufio.Reader
}

func newRotationService(option Options) (r Rotation, err error) {
	r = Rotation{}
	r.d, err = newAdbDevice(option.Serial, option.Adb)
	r.closed = true
	return
}

//download file to device
func (r *Rotation) download(path, url string) (err error) {
	fout, err := r.d.Device.OpenWrite(path, 0755, time.Now())
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

//install rotationWatcher.apk
func (r *Rotation) install() (err error) {
	url := "https://github.com/NetEaseGame/AutomatorX/raw/master/atx/vendor/RotationWatcher.apk"
	//check package
	pkgName := "jp.co.cyberagent.stf.rotationwatcher"
	plist, err := r.d.getPackageList()
	if err != nil {
		return
	}
	for _, val := range plist {
		if strings.Contains(val, pkgName) {
			return
		}
	}
	//downlaod apk
	path := "/data/local/tmp/RotationWatcher.apk"
	err = r.download(path, url)
	if err != nil {
		return
	}
	//install apk
	_, err = r.d.shell("pm", "install", "-rt", path)
	if err != nil {
		return
	}
	return
}

//start rotation service
func (r *Rotation) start() (err error) {
	pkgName := "jp.co.cyberagent.stf.rotationwatcher"
	out, err := r.d.shell("pm path " + pkgName)
	if err != nil {
		return
	}
	fields := strings.Split(strip(out), ":")
	path := fields[len(fields)-1]
	r.proc = r.d.buildCommand("CLASSPATH="+path, "app_process", "/system/bin", "jp.co.cyberagent.stf.rotationwatcher.RotationWatcher")
	r.proc.Stderr = os.Stderr
	stdoutReader, err := r.proc.StdoutPipe()
	if err != nil {
		return
	}
	r.brd = bufio.NewReader(stdoutReader)
	return r.proc.Start()
}

func (r *Rotation) watch() (orienC <-chan int, err error) {
	rC := make(chan int, 0)
	go func() {
		for {
			line, _, er := r.brd.ReadLine()
			if er != nil {
				// restart if rotation watcher crashed.
				r.start()
				continue
				//break
			}
			tmp := strings.Replace(string(line), "\r", "", -1)
			tmp = strings.Replace(tmp, "\n", "", -1)
			orientation, er := strconv.Atoi(string(tmp))
			if er != nil {
				break
			}
			rC <- orientation
		}
		close(rC)
	}()
	orienC = rC
	return
}
