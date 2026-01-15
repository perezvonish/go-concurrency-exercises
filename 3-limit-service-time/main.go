//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	sync.RWMutex
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) IncrementUsedTimeBySecond() {
	u.Lock()
	defer u.Unlock()

	u.TimeUsed = u.TimeUsed + 1
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if !isLimitAvailable(u) {
		return false
	}

	done := make(chan struct{}, 1)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		process()
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			return true
		case <-ticker.C:
			u.IncrementUsedTimeBySecond()
			if !isLimitAvailable(u) {
				return false
			}
		}
	}
}

func isLimitAvailable(u *User) bool {
	if !u.IsPremium && u.TimeUsed >= 10 {
		return false
	}

	return true
}

func main() {
	RunMockServer()
}
