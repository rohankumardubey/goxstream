package main

import (
	"fmt"
	"goxstream/internal/api"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	killProcess()
	fmt.Println("GoXStream REST API running on :8080")
	if err := api.StartAPIServer(":8080"); err != nil {
		panic(err)
	}
}

func killProcess() {
	out, err := exec.Command("lsof", "-t", "-i:8080").Output()
	if err != nil || len(out) == 0 {
		return // Nothing running or lsof not available
	}
	pids := strings.Fields(string(out))
	for _, pidStr := range pids {
		pid, _ := strconv.Atoi(pidStr)
		proc, _ := os.FindProcess(pid)
		proc.Signal(syscall.SIGKILL)
	}
}
