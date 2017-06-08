package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/util"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
var input = flag.String("input", "auto", "input (auto / adb / stdin / <filename>)")

var highlightPattern, ftagPattern string
var termcols int
var outputFile *os.File
var pids []int

func main() {
	termcols = getTermWidth()

	flag.Parse()
	buildPatterns()

	if len(*file) > 0 {
		f, err := os.Create(*file)
		if err != nil {
			log.Fatal("Error opening output file: "+*file, err)
		}
		outputFile = f
	}

	if *input == "auto" {
		if !util.IsTerminal(os.Stdin) {
			*input = "stdin"
		} else {
			*input = "adb"
		}
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

func buildPatterns() {
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
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.Replace(pattern, "\\*", ".*", -1)
	pattern = "^" + pattern + "$"
	if !*casesensitive {
		pattern = "(?i)" + pattern
	}

	return pattern
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
		if termcols > 0 {
			message = wrapmessage(message)
		}
	} else if *color {
		// Apply color (based on priority) otherwise
		pre = colorMap[prio]
	}

	// Limit tag (if necessary)
	if len(tag) > *taglength {
		tag = tag[0:*taglength]
	}

	// Print to stdout
	if *stdout {
		fmt.Printf("%s%-"+strconv.Itoa(*taglength)+"s[%s] %s%s\n", pre, tag, prio, message, Reset)
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
	m, _ := regexp.MatchString(pattern, strings.TrimSpace(s))
	return m
}
