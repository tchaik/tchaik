// +build !unix

package main

import "time"

func getFileCreationTime(string) (time.Time) {
	return time.Time{}
}
