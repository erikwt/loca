package main

const (
	DEFAULT_TAG_LENGTH  = 30
	DEFAULT_PRIO_FILTER = "VDIWEF"
	DEFAULT_MINPRIO     = "V"

	REGEXP_ADB_STD        = "(?P<prio>.)/(?P<tag>.+)\\(\\s*\\d+\\):\\s+(?P<msg>.+)"
	REGEXP_ADB_THREADTIME = "\\d+-\\d+\\s+\\d+:\\d+:\\d+.\\d+\\s+\\d+\\s+\\d+\\s+(P<prio>.)\\s+(P<tag>.+):\\s+(P<msg>.+)"
	REGEXP_PUSSLOG_STD    = "(?P<tag>.+)\\s*\\[(?P<prio>.)\\]\\s+(?P<msg>.+)"
)

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
