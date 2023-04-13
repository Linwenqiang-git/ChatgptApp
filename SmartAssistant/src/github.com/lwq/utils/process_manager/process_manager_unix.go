package processmanager

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type ProcessManagerLinux struct {
}

func (p *ProcessManagerLinux) Kill(port int) error {
	log.Println("come in win process kill...")
	cmdStr := fmt.Sprintf("lsof -i:%s", strconv.Itoa(port))
	cmd := exec.Command("cmd", "/C", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// 解析进程PID
	pid := ""
	// linux下的PID信息格式为“COMMAND   PID     USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME”
	fields := strings.Fields(string(output))
	if len(fields) >= 2 {
		pid = fields[1]
	}
	// 关闭进程
	if pid != "" {
		fmt.Printf("Closing process %s on port %s...\n", pid, strconv.Itoa(port))
		cmd := exec.Command("taskkill", "/F", "/PID", pid)
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Println("Process closed successfully.")
	} else {
		fmt.Printf("No process found on port %s.\n", strconv.Itoa(port))
	}
	return nil
}
