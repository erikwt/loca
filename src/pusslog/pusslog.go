package main

import (
	"log"
	"os/exec"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"flag"
)

var process = flag.String("p", "", "process or package name filter") // TODO
var highlight = flag.String("hl", "", "highlight tag/process/package name") // TODO (tag implemented)
var prio = flag.String("prio", "VDIWEF", "priority filter (VERBOSE/DEBUG/INFO/WARNING/ERROR/FATAL)") // TODO
var minprio = flag.String("minprio", "V", "minimum priority level") // TODO

func main() {
    flag.Parse()
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
    var pre, post string
    if tag == *highlight {
        pre = BgYellow + FgBlack
        post = Reset
    }
    fmt.Printf("%s[%s] %s%s\n", pre, tag, message, post)
}

