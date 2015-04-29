// +build unix

package main

import (
	"time"

	"golang.org/x/sys/unix"
)

func getFileCreationTime(path string) (time.Time) {
	stat := unix.Stat_t{}
	unix.Lstat(path, &stat)

	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}
