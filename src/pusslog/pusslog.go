package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	DEFAULT_TAG_LENGTH  = 30
	DEFAULT_PRIO_FILTER = "VDIWEF"
	DEFAULT_MINPRIO     = "V"

	REGEXP_ADB_STD        = "(?P<prio>.)/(?P<tag>.+)\\(\\s*\\d+\\):\\s+(?P<msg>.+)"
	REGEXP_ADB_THREADTIME = "\\d+-\\d+\\s+\\d+:\\d+:\\d+.\\d+\\s+\\d+\\s+\\d+\\s+(P<prio>.)\\s+(P<tag>.+):\\s+(P<msg>.+)"
	REGEXP_PUSSLOG_STD    = "(?P<tag>.+)\\s*\\[(?P<prio>.)\\]\\s+(?P<msg>.+)"
)

var taglength = flag.Int("tl", DEFAULT_TAG_LENGTH, "maximum tag length")
var process = flag.String("p", "", "process or package name filter")
var ftag = flag.String("t", "", "tag filter")
var highlight = flag.String("hl", "", "highlight tag/process/package name")
var priofilter = flag.String("prio", DEFAULT_PRIO_FILTER, "priority filter (VERBOSE/DEBUG/INFO/WARNING/ERROR/FATAL)")
var minprio = flag.String("minprio", DEFAULT_MINPRIO, "minimum priority level")
var file = flag.String("file", "", "write log to file")
var grep = flag.String("grep", "", "grep on log message (regex filter)")
var color = flag.Bool("color", true, "enable colored output")
var stdout = flag.Bool("stdout", true, "print to <stdout>")
var casesensitive = flag.Bool("casesensitive", false, "case sensitive filters")
var input = flag.String("input", "adb", "input (adb / stdin / <filename>")

var prioMap = map[string]int{
	"V": 0,
	"D": 1,
	"I": 2,
	"W": 3,
	"E": 4,
	"F": 5,
}

var colorMap = map[string]string{
	"V": FgGreen,
	"D": FgCyan,
	"I": FgYellow,
	"W": FgBlue,
	"E": FgRed,
	"F": FgMagenta,
}

var highlightMap = map[string]string{
	"V": BgGreen + FgBlack,
	"D": BgCyan + FgBlack,
	"I": BgYellow + FgBlack,
	"W": BgBlue + FgBlack,
	"E": BgRed + FgBlack,
	"F": BgMagenta + FgBlack,
}

var processPattern, highlightPattern, ftagPattern string
var termcols int
var outputFile *os.File
var pids []int

func main() {
	termcols = GetWinsize()

	flag.Parse()
	buildPatterns()

	if len(*file) > 0 {
		f, err := os.Create(*file)
		if err != nil {
			log.Fatal("Error opening output file: "+*file, err)
		}
		outputFile = f
	}

	switch *input {
	case "adb":
		testEnv()
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

		getPids()

		adbReadlog(deviceId)

	case "stdin":
		fileReadlog(os.Stdin)

	default:
		file, err := os.Open(*input)
		if err != nil {
			log.Fatal("Error: ", err)
			return
		}

		fileReadlog(file)
	}

}

func testEnv() {
	if _, err := exec.LookPath("adb"); err != nil {
		log.Fatal("Error: adb command not found in PATH")
	}
}

func buildPatterns() {
	if len(*process) > 0 {
		processPattern = buildPattern(*process)
	}

	if len(*highlight) > 0 {
		highlightPattern = buildPattern(*highlight)
	}

	if len(*ftag) > 0 {
		ftagPattern = buildPattern(*ftag)
	}

	if !*casesensitive && len(*grep) > 0 {
		*grep = "(?i)" + *grep
	}
}

func buildPattern(pattern string) string {
	if !*casesensitive {
		pattern = strings.ToLower(pattern)
	}

	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.Replace(pattern, "\\*", ".*", -1)
	return "^" + pattern + "$"
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

func fileReadlog(file *os.File) {
	rd := bufio.NewReader(file)

	formatAdbStd, err := regexp.Compile(REGEXP_ADB_STD)
	if err != nil {
		log.Fatal("Regexp compile error", err)
	}
	formatAdbThreadtime, err := regexp.Compile(REGEXP_ADB_THREADTIME)
	if err != nil {
		log.Fatal("Regexp compile error", err)
	}
	formatPusslogStd, err := regexp.Compile(REGEXP_PUSSLOG_STD)
	if err != nil {
		log.Fatal("Regexp compile error", err)
	}

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal("Read Error: ", err)
			}
			return
		}

		var format *regexp.Regexp
		if formatAdbStd.MatchString(str) {
			format = formatAdbStd
		} else if formatAdbThreadtime.MatchString(str) {
			format = formatAdbThreadtime
		} else if formatPusslogStd.MatchString(str) {
			format = formatPusslogStd
		} else {
			log.Println("Does not match any format: " + str)
			continue
		}

		var tag, prio, msg string
		matches := format.FindStringSubmatch(str)
		names := format.SubexpNames()
		for i, n := range names {
			if n == "tag" {
				tag = matches[i]
			} else if n == "prio" {
				prio = matches[i]
			} else if n == "msg" {
				msg = matches[i]
			}
		}

		logmessage("", "", 0, 0, prio, tag, msg)
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

func logmessage(date string, time string, threadid int, processid int, prio string, tag string, message string) {
	// process id filter (if enabled)
	if len(*process) > 0 && !contains(pids, processid) {
		return
	}

	// Tag filter (if enabled)
	if len(*ftag) > 0 && !matches(tag, ftagPattern) {
		return
	}

	// prio filter
	if !strings.Contains(*priofilter, prio) {
		return
	}

	// min prio filter
	if prioMap[*minprio] > prioMap[prio] {
		return
	}

	// grep filter
	if len(*grep) > 0 {
		if matches, _ := regexp.MatchString(*grep, message); !matches {
			return
		}
	}

	// highlight (if enabled)
	var pre string
	if (len(*highlight) > 0 && matches(tag, highlightPattern)) || (len(*process) == 0 && contains(pids, processid)) {
		pre = highlightMap[prio]
	} else if *color {
		// Apply color (based on priority) otherwise
		pre = colorMap[prio]
	}

	// Wrap message and fill to terminal width
	wrappedmessage := wrapmessage(message)

	// Limit tag (if necessary)
	if len(tag) > *taglength {
		tag = tag[0:*taglength]
	}

	// Print to stdout
	if *stdout {
		fmt.Printf("%s%-"+strconv.Itoa(*taglength)+"s[%s] %s%s\n", pre, tag, prio, wrappedmessage, Reset)
	}

	// Print to file if needed
	if len(*file) > 0 {
		message = fmt.Sprintf("%-"+strconv.Itoa(*taglength)+"s[%s] %s\n", tag, prio, message)
		if _, err := outputFile.Write([]byte(message)); err != nil {
			log.Fatal("Error writing to logfile: "+*file, err)
		}
	}
}

func wrapmessage(message string) string {
	if termcols == -1 {
		return message
	}
	
	availableWidth := termcols - *taglength - 4
	parts := len(message) / availableWidth
	if len(message)%availableWidth != 0 {
		parts++
	}
	if parts > 1 {
		var newmessage string
		var end int
		start := 0
		for {
			end = start + availableWidth
			if end > len(message) {
				end = len(message)
			}
			newmessage += message[start:end]

			if end < len(message) {
				start = end
				newmessage += "\n"
				for i := 0; i < *taglength+4; i++ {
					newmessage += " "
				}
			} else {
				numSpaces := availableWidth - (end - start)
				for i := 0; i < numSpaces; i++ {
					newmessage += " "
				}

				break
			}
		}
		message = newmessage
	} else {
		numSpaces := availableWidth - len(message)
		for i := 0; i < numSpaces; i++ {
			message += " "
		}
	}

	return message
}

func contains(list []int, elem int) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func matches(s string, pattern string) bool {
	if !*casesensitive {
		s = strings.ToLower(s)
	}

	m, _ := regexp.MatchString(pattern, s)
	return m
}
