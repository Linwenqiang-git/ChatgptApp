package processmanager

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type ProcessManagerWin struct {
}

func (p *ProcessManagerWin) Kill(port int) error {
	log.Println("come in win process kill...")
	cmdStr := fmt.Sprintf("netstat -ano | findstr :%s", strconv.Itoa(port))
	cmd := exec.Command("cmd", "/C", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	// 解析进程PID
	pid := ""
	// windows下的PID信息格式为“TCP    0.0.0.0:8080    0.0.0.0:0    LISTENING    PID”
	for _, str := range strings.Split(string(output), "\n") {
		fields := strings.Fields(str)
		if len(fields) >= 5 && fields[3] == "LISTENING" {
			pid = fields[len(fields)-1]
			break
		}
	}
	// 关闭进程
	if pid != "" {
		cmd := exec.Command("taskkill", "/F", "/PID", pid)
		if err := cmd.Run(); err != nil {
			return err
		}
		log.Printf("Closing process %s on port %s...\n", pid, strconv.Itoa(port))
	} else {
		log.Printf("No process found on port %s.\n", strconv.Itoa(port))
	}
	return nil
}
