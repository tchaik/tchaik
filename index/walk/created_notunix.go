// +build !unix

package walk

import "time"

func getCreatedTime(string) (time.Time, error) {
	return time.Time{}, nil
}
