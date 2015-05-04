// +build !unix

package main

import "time"

func getCreatedTime(string) (time.Time, error) {
	return time.Time{}, nil
}
