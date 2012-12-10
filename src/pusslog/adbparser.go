package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func testEnv() {
	if _, err := exec.LookPath("adb"); err != nil {
		log.Fatal("Error: adb command not found in PATH")
	}
}

func getDeviceId() (string, error) {
	cmd := exec.Command("adb", "devices")
	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error getting devices: %s", err)
	}

	// Skip irrelevant lines
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			return "", errors.New("Error getting devices")
		}
		if len(str) > 0 && strings.TrimSpace(str)[0] != '*' {
			break
		}
	}

	devices := make([]string, 0)
	for str, err := rd.ReadString('\n'); err == nil; str, err = rd.ReadString('\n') {
		if str = strings.TrimSpace(str); len(str) > 0 {
			devices = append(devices, str)
		}
	}

	if len(devices) == 0 {
		return "", errors.New("No device connected")
	}

	if len(devices) == 1 {
		f := strings.Fields(devices[0])
		return f[0], nil
	}

	fmt.Println("Multiple devices found!\n")
	for i := 0; i < len(devices); i++ {
		fmt.Printf("[%d]\t%s\n", i+1, devices[i])
	}

	deviceIndex := 0
	for deviceIndex <= 0 || deviceIndex > len(devices) {
		fmt.Printf("\nUse device number: ")
		fmt.Scanf("%d", &deviceIndex)
	}

	return strings.Fields(devices[deviceIndex-1])[0], nil
}

func getPids() {
	pids = make([]int, 0)

	if len(*process) > 0 {
		addPids(*process)
	}

	if len(*highlight) > 0 {
		addPids(*highlight)
	}
}

func addPids(processname string) {
	cmd := exec.Command("adb", "shell", "ps")

	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	// Skip first line
	if _, err := rd.ReadString('\n'); err != nil {
		return
	}

	for str, err := rd.ReadString('\n'); err == nil; str, err = rd.ReadString('\n') {
		if fields := strings.Fields(str); len(fields) == 9 && matches(fields[8], processname) {
			pid, _ := strconv.Atoi(fields[1])
			pids = append(pids, pid)
		}
	}
}

func adbReadlog(deviceId string) {
	cmd := exec.Command("adb", "-s", deviceId, "logcat", "-v", "threadtime")

	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Device disconnected.")
			} else {
				log.Fatal("Read Error: ", err)
			}
			return
		}

		parseline(str)
	}
}

func parseline(l string) {
	fields := strings.Fields(l)
	if len(fields) >= 7 {
		date := fields[0]
		time := fields[1]
		threadid, _ := strconv.Atoi(fields[2])
		processid, _ := strconv.Atoi(fields[3])
		prio := fields[4]
		tag := strings.TrimRight(fields[5], ":")
		message := strings.TrimLeft(strings.Join(fields[6:], " "), ": ")

		logmessage(date, time, threadid, processid, prio, tag, message)

		if "ActivityManager" == tag &&
			(len(*process) > 0 && strings.Contains(message, *process) ||
				len(*highlight) > 0 && strings.Contains(message, *highlight)) {

			getPids()
		}
	}
}
