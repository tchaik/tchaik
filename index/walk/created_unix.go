// +build unix

package walk

import (
	"time"

	"golang.org/x/sys/unix"
)

func getCreatedTime(path string) (time.Time, error) {
	stat := unix.Stat_t{}
	err := unix.Lstat(path, &stat)
	if err != nil {
		return err
	}
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}
