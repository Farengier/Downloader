package src

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	max := 30000
	for i := 29900; i <= max; i++ {
		_renderProgress(max, i)
		time.Sleep(time.Millisecond)
	}
}
