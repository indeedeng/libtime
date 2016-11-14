package libtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Clock_Now(t *testing.T) {
	standardNow := time.Now()

	c := SystemClock()
	fake := c.Now()

	require.True(t, fake.After(standardNow))
}

func Test_Clock_Since(t *testing.T) {
	standardNow := time.Now()

	c := SystemClock()

	duration := c.Since(standardNow)
	require.True(t, duration > 0)
}

func Test_Clock_SinceMS(t *testing.T) {
	standardNow := time.Now()
	time.Sleep(2 * time.Millisecond)

	c := SystemClock()

	duration := c.SinceMS(standardNow)
	require.True(t, duration > 0)
}
