package settings

import (
	"errors"
	"runtime"
	"strconv"

	"github.com/spf13/viper"
)

type IpcServerSetting struct {
	addr string
	port int
	fd   string
}

func (s *IpcServerSetting) GetConnectInfo() (protocol string, address string) {
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
	return 0, errors.New("no listening port is required")
}

func NewIpcServerSetting() *IpcServerSetting {
	return &IpcServerSetting{
		addr: viper.GetString("ipc_server.addr"),
		port: viper.GetInt("ipc_server.port"),
		fd:   viper.GetString("ipc_server.fd"),
	}
}
