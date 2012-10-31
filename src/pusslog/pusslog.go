package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var process = flag.String("p", "", "process or package name filter")                                 // TODO
var highlight = flag.String("hl", "", "highlight tag/process/package name")                          // TODO (tag implemented)
var prio = flag.String("prio", "VDIWEF", "priority filter (VERBOSE/DEBUG/INFO/WARNING/ERROR/FATAL)") // TODO
var minprio = flag.String("minprio", "V", "minimum priority level")                                  // TODO

var pid int

func main() {
	flag.Parse()
	if len(*process) > 0 || len(*highlight) > 0 {
		setPidFilter()
	}

	loop()
}

func loop() {
	cmd := exec.Command("adb", "logcat", "-v", "threadtime")

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
	if pid > 0 && pid != processid {
		return
	}

	// highlight (if enabled)
	var pre, post string
	if tag == *highlight {
		pre = BgYellow + FgBlack
		post = Reset
	}

	// Print logmessage
	fmt.Printf("%s[%s] %s%s\n", pre, tag, message, post)
}

func setPidFilter() {
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
