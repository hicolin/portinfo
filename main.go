package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var usage = `Usage: portinfo <port>`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	port := os.Args[1]
	_, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		fmt.Println(usage)
		return
	}

	// 使用 netstat 获取端口信息
	cmd := exec.Command("cmd", "/c", "netstat -ano | findstr :"+port+" ")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			fmt.Printf("Port %s not found\n", port)
			return
		}
		fmt.Printf("Error running netstat: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%s ", port)) {
			parts := strings.Fields(line)
			if len(parts) < 5 {
				continue
			}
			pid := parts[len(parts)-1]
			fmt.Printf("Port %s is used by PID %s\n", port, pid)

			// 使用 tasklist 获取进程信息
			taskCmd := exec.Command("tasklist", "/FI", "PID eq "+pid)
			taskOutput, err := taskCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running tasklist: %v\n", err)
				return
			}

			taskLines := strings.Split(string(taskOutput), "\n")
			for _, taskLine := range taskLines {
				taskParts := strings.Fields(taskLine)
				if len(taskParts) > 0 && strings.Contains(taskParts[0], ".exe") {
					fmt.Printf("Process Name: %s\n", taskParts[0])
				}
			}
			break
		}
	}
}
