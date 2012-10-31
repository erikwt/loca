package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var process = flag.String("p", "", "process or package name filter")
var highlight = flag.String("hl", "", "highlight tag/process/package name")
var priofilter = flag.String("prio", "VDIWEF", "priority filter (VERBOSE/DEBUG/INFO/WARNING/ERROR/FATAL)")
var minprio = flag.String("minprio", "V", "minimum priority level")

var pid int

func main() {
	flag.Parse()

	deviceId, err := getDeviceId()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	if deviceId == "????????????" {
		log.Fatal("No permissions for device")
		return
	}

	fmt.Printf("Selected device: %s\n\n", deviceId)

	if len(*process) > 0 || len(*highlight) > 0 {
		setPid()
	}

	loop(deviceId)
}

func getDeviceId() (string, error) {
	cmd := exec.Command("adb", "devices")
	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Error getting devices: %s", err)
	}

	// Skip first line
	if _, err := rd.ReadString('\n'); err != nil {
		return "", errors.New("Error getting devices")
	}

	devices := make([]string, 0)
	for str, err := rd.ReadString('\n'); err == nil; str, err = rd.ReadString('\n') {
		if len(strings.TrimSpace(str)) > 0 {
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
		fmt.Printf("[%d]\t%s", i+1, devices[i])
	}

	deviceIndex := 0
	for deviceIndex <= 0 || deviceIndex > len(devices) {
		fmt.Printf("\nUse device number: ")
		fmt.Scanf("%d", &deviceIndex)
	}

	return strings.Fields(devices[deviceIndex-1])[0], nil
}

func setPid() {
	var cmd = new(exec.Cmd)
	if len(*process) > 0 {
		cmd = exec.Command("adb", "shell", "ps", *process)
	} else if len(*highlight) > 0 {
		cmd = exec.Command("adb", "shell", "ps", *highlight)
	}

	if cmd == nil {
		return
	}

	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	// Skip first line
	_, err := rd.ReadString('\n')
	if err != nil {
		log.Fatal("Read Error:", err)
		return
	}

	str, err := rd.ReadString('\n')
	if err != nil {
		log.Fatal("Read Error:", err)
		return
	}

	fields := strings.Fields(str)
	if len(fields) == 9 {
		pid, _ = strconv.Atoi(fields[1])
	}
}

func loop(deviceId string) {
	cmd := exec.Command("adb", "-s", deviceId, "logcat", "-v", "threadtime")

	stdout, _ := cmd.StdoutPipe()
	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
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
	}
}

func logmessage(date string, time string, threadid int, processid int, prio string, tag string, message string) {
	// process id filter (if enabled)
	if len(*process) > 0 && pid != processid {
		return
	}
	
	// prio filter
	if !strings.Contains(*priofilter, prio) {
	    return
    }

	// highlight (if enabled)
	var pre, post string
	if tag == *highlight || (len(*process) == 0 && pid == processid) {
		pre = BgYellow + FgBlack
		post = Reset
	}

	// Print logmessage
	fmt.Printf("%s[%s] %s%s\n", pre, tag, message, post)
}
