// +build windows

package git

import (
	"syscall"
	"time"

	"github.com/newsletter2go/go-git/plumbing/format/index"
)

func init() {
	fillSystemInfo = func(e *index.Entry, sys interface{}) {
		if os, ok := sys.(*syscall.Win32FileAttributeData); ok {
			seconds := os.CreationTime.Nanoseconds() / 1000000000
			nanoseconds := os.CreationTime.Nanoseconds() - seconds*1000000000
			e.CreatedAt = time.Unix(seconds, nanoseconds)
		}
	}
}
