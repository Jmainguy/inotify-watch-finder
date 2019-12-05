# inotify-watch-finder
This tool is designed to help find what files / directories are being watched by a specific PID.

We had 65,000 inotify's being consumed by the kubelet on one of our hosts, and used this tool to see what it was watching (/sys/fs/cgroup it turns out)

## Usage
1. First, find out which PID you want to check
2. Find the inotify files in fdinfo
```/bin/bash
ls -l /proc/<PID>/fd | grep -i inotify
```
3. Start catting the same number you found above in /fdinfo, and find the file with tons of lines in it instead of 4-5
```/bin/bash
cat /proc/<PID>/fdinfo/26
```
4. Use inotify-watch-finder to track down the files being watched.

This tool finds the inode for each watch, and then it finds the inode of every file on the system, / specifically.

Then it finds the filename for each inode in the first list and prints it out, giving you the name of all files being watched by that process.

This takes along time, on our system, about 40 minutes. I suggest piping it to a text file and reviewing it later at your convenience.

```/bin/bash
nice ./inotify-watch-finder -fdinfo /proc/123031/fdinfo/23 > /tmp/inotifyReport
```

## Build
```/bin/bash
export GO111MODULE=on
go mod init
go build
```

