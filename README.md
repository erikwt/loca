# Pusslog

A commandline Android 'logcat' reader with many additional features.
Written in Go, for the purpose of learning Go and improving the Android development experience.

## Features
* Colored output
* Write to file
* Highlight messages from package (like com.example)
* Highlight messages with specific tag
* Filter messages from package (like com.example)
* Filter specific priorities
* Filter by minimum priority

## Building
Pusslog is written in Go. To build pusslog, you'll need to install go on your system.

Build pusslog:
``` bash
export GOPATH=/path/to/pusslog/source
go build pusslog
```

Now you can run the 'pusslog' executable:
``` bash
./pusslog
```

Or put the executable somewhere in your $PATH.

## License
Pusslog is licenced under the GNU GPL version 3.

## Author
I'm a young, enthusiastic hacker from Amsterdam. I study computer science at the VU (Free University) and my work mostly involves Android programming. I do a lot of hacking in my spare time, resulting in many projects I want to share with the world.

* Twitter: [@ewallentinsen](http://www.twitter.com/ewallentinsen)
* Blog: http://www.erikw.eu
