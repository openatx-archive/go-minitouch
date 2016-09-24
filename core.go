/*
* File: core.go
* Author : bigwavelet
* Description: core file
* Created: 2016-08-26
 */

package minitouch

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"unicode"
)

func strip(str string) string {
	return strings.TrimFunc(str, unicode.IsSpace)
}

func splitLines(str string) (result []string) {
	tmp := strings.Replace(str, "\r\n", "\n", -1)
	tmp = strings.Replace(str, "\r", "", -1)
	result = strings.Split(tmp, "\n")
	return
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randPort() (port int, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	addr := listener.Addr().String()
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(portString)
}
