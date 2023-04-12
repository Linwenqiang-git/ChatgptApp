package settings

import (
	"errors"
	"runtime"
	"strconv"

	"gopkg.in/ini.v1"
)

type IpcServerSetting struct {
	addr string
	port int
	fd   string
}

func (s *IpcServerSetting) GetConnetInfo() (protocol string, address string) {
	sysType := runtime.GOOS
	if sysType == "windows" {
		protocol = "tcp"
		address = s.addr + ":" + strconv.Itoa(s.port)
	} else {
		protocol = "unix"
		address = s.fd
	}
	return
}
func (s *IpcServerSetting) GetPort() (int, error) {
	sysType := runtime.GOOS
	if sysType == "windows" {
		return s.port, nil
	}
	return 0, errors.New("No listening port is required.")
}

func NewIpcServerSetting(section *ini.Section) *IpcServerSetting {
	return &IpcServerSetting{
		addr: section.Key("addr").String(),
		port: section.Key("port").MustInt(),
		fd:   section.Key("fd").String(),
	}
}
