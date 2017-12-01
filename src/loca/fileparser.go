package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
)

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
	formatLocaStd, err := regexp.Compile(REGEXP_LOCA_STD)
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
		} else if formatLocaStd.MatchString(str) {
			format = formatLocaStd
		} else {
			log.Print("Does not match any known logformat: " + str)
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
