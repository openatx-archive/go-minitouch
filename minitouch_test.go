/*
* File: minicap_test.go
* Author : bigwavelet
* Description: android minicap test
* Created: 2016-09-13
 */

package minitouch

import "testing"

var option = Options{"EP7333W7XB", 0, ""}

func TestNewService(t *testing.T) {
	_, err := NewService(option)
	if err != nil {
		t.Error("New Service error:" + err.Error())
	} else {
		t.Log("New Service Test Passed.")
	}
}

func TestInstall(t *testing.T) {
	m, err := NewService(option)
	if err != nil {
		t.Error("New Service error:" + err.Error())
		return
	}
	err = m.Install()
	if err != nil {
		t.Error("minitouch service Install Test error:" + err.Error())
	} else {
		t.Log("Install Test Passed.")
	}
}

func TestIsSupported(t *testing.T) {
	m, err := NewService(option)
	if err != nil {
		t.Error("New Service error:" + err.Error())
		return
	}
	err = m.Install()
	if err != nil {
		t.Error("minitouch service Install Test error:" + err.Error())
		return
	}
	supported := m.IsSupported()
	if supported {
		t.Log("minitouch supported.")
	} else {
		t.Log("minitouch not supported.")
	}
}
