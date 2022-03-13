package libtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_SafeTimer(t *testing.T) {
	delay := 3 * time.Millisecond
	start := time.Now()

	timer, stop := SafeTimer(delay)
	defer stop()
	<-timer.C

	elapsed := time.Since(start)
	require.GreaterOrEqual(t, int64(elapsed), int64(delay))
}
