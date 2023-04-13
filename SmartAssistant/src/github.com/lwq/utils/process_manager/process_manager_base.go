package processmanager

import "runtime"

type IProcessManager interface {
	Kill(port int) error
}

func NewProcesManager() IProcessManager {
	if runtime.GOOS == "windows" {
		return &ProcessManagerWin{}
	} else {
		return &ProcessManagerLinux{}
	}
}
