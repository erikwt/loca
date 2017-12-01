# LOCA

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
There is a Linux x86-64 binary available in [the download section](https://github.com/erikwt/loca/downloads/).
To use loca on other platforms, just build the binary from source as described below.

## Building
Loca is written in Go. To build loca, you'll need to install go on your system.

Clone loca from github:
``` bash
git clone https://github.com/erikwt/loca.git
```

Build loca:
``` bash
cd /path/to/loca/source # path where you cloned the source
export GOPATH=$(pwd)
go get github.com/pebbe/util
go build loca
```

Now you can run the 'loca' executable:
``` bash
./loca
```

Or put the executable somewhere in your $PATH.

## Example usages
Note that many options can be combined to allow for very complex filtering and highlighting. Loca accepts
input directly from a device (over adb), from a file or from the standard input (pipe).

``` bash
# Same as 'adb logcat' but prettier
loca

# loca help (brief) with all the options and defaults
loca -help

# highlights all log messages from the app with package 'eu.erikw.myapp' 
loca -hl eu.erikw.myapp

# only log messages from the app with package 'eu.erikw.myapp'
loca -p eu.erikw.myapp

# only messages with priority 'Error' or 'Fatal'
loca -prio EF

# only messages with priority 'Warning' or higher
loca -minprio W

# only messages with tag 'MyAppTag'
loca -t MyAppTag

# write log to file
loca -file /tmp/logfile

# read from logfile instead of device
loca -input /tmp/logfile

# read from stdin (adb logcat output without messages containing 'dalvikvm' 
adb logcat | grep -v dalvikvm | loca

# only messages matching regexp 'http://.+.com'
loca -grep http://.+.com

# combination of filters and use of wildcards ;-)
loca -t *man* -hl eu.erikw.myapp -minprio I
```

## Problems / issues / bugs
Please report issues in the [issue tracker](https://github.com/erikwt/loca/issues).

## License
Loca is licenced under the GNU GPL version 3.

## Open source
Thanks to [pebbe](https://github.com/pebbe) for the IsTerminal(..) implementation in the [go util source](https://github.com/pebbe/util).

## Author
I'm a young, enthusiastic hacker from Amsterdam. I study computer science at the VU (Free University) and my work mostly involves Android programming. I do a lot of hacking in my spare time, resulting in many projects I want to share with the world.

* Twitter: [@ewallentinsen](http://www.twitter.com/ewallentinsen)
* Blog: http://www.erikw.eu
