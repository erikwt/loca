# Pusslog

A commandline Android 'logcat' reader with many additional features.
Written in Go, for the purpose of learning Go and improving the Android development experience.

## Features
* Colored output
* Write to file
* Highlight messages from package (like com.example)
* Highlight messages with specific tag
* Filter messages from package (like com.example)
* Filter messages from specific tag
* Filter specific priorities
* Filter by minimum priority
* Filter by regexp in log message (grep)
* Wildcard and case insensitive filtering
* Read logs from device (adb), logfile or standard input (stdin)

Example usages below.

## Download
There is a Linux x86-64 binary available in [the download section](https://github.com/erikwt/pusslog/downloads/).
To use pusslog on other platforms, just build the binary from source as described below.

## Building
Pusslog is written in Go. To build pusslog, you'll need to install go on your system.

Clone pusslog from github:
``` bash
git clone https://github.com/erikwt/pusslog.git
```

Build pusslog:
``` bash
cd /path/to/pusslog/source # path where you cloned the source
export GOPATH=$(pwd)
go build pusslog
```

Now you can run the 'pusslog' executable:
``` bash
./pusslog
```

Or put the executable somewhere in your $PATH.

## Example usages
``` bash
# Same as 'adb logcat' but prettier
pusslog

# pusslog help (brief) with all the options and defaults
pusslog -help

# highlights all log messages from the app with package 'eu.erikw.myapp' 
pusslog -hl eu.erikw.myapp

# only log messages from the app with package 'eu.erikw.myapp'
pusslog -p eu.erikw.myapp

# only messages with priority 'Error' or 'Fatal'
pusslog -prio EF

# only messages with priority 'Warning' or higher
pusslog -minprio W

# only messages with tag 'MyAppTag'
pusslog -t MyAppTag

# write log to file
pusslog -file /tmp/logfile

# read from logfile instead of device
pusslog -input /tmp/logfile

# read from stdin (adb logcat output without messages containing 'dalvikvm' 
adb logcat | grep -v dalvikvm | pusslog

# only messages matching regexp 'http://.+.com'
pusslog -grep http://.+.com

# combination of filters and use of wildcards ;-)
pusslog -t *man* -hl eu.erikw.myapp -minprio I
```

## License
Pusslog is licenced under the GNU GPL version 3.

## Author
I'm a young, enthusiastic hacker from Amsterdam. I study computer science at the VU (Free University) and my work mostly involves Android programming. I do a lot of hacking in my spare time, resulting in many projects I want to share with the world.

* Twitter: [@ewallentinsen](http://www.twitter.com/ewallentinsen)
* Blog: http://www.erikw.eu
