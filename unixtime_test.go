package libtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ToMilliseconds(t *testing.T) {
	seconds := int64(100)
	nanoseconds := int64(1000000)

	expected := seconds*int64(time.Second) + nanoseconds*int64(time.Nanosecond)
	actual := ToMilliseconds(time.Unix(seconds, nanoseconds)) * int64(time.Millisecond)
	require.Equal(t, expected, actual)
}

func Test_FromMilliseconds(t *testing.T) {
	milliseconds := int64(1120)

	expected := time.Unix(0, 0).Add(time.Duration(milliseconds) * time.Millisecond)
	actual := FromMilliseconds(milliseconds)

	require.Equal(t, expected, actual)
}

func Test_DurationToMillis(t *testing.T) {
	duration := 1234567890 * time.Nanosecond
	require.Equal(t, int64(1234), DurationToMillis(duration))
}
