package libtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Sleeper(t *testing.T) {
	s := NewSleeper()
	before := time.Now()
	s.Sleep(2 * time.Millisecond)
	after := time.Now()

	diff := after.Sub(before)
	require.True(t, diff.Nanoseconds() >= int64(2*time.Millisecond))
}
