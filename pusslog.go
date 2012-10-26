package main

import (
	"log"
	"os/exec"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

func main() {
    // TODO: Implement options
	loop()
}

func loop(){
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

func parseline(l string){
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

func logmessage(date string, time string, threadid int, processid int, prio string, tag string, message string){
    fmt.Printf("[%s] %s\n", tag, message)
}

