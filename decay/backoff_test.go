package decay

import (
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libtime"
)

var (
	errOops = errors.New("oops")
)

func echo(keepRetrying bool, err error) TryFunc {
	return func() (bool, error) {
		return keepRetrying, err
	}
}

type iteration struct {
	keepTrying bool
	err        error
}

func schedule(iterations []iteration) TryFunc {
	i := 0
	return func() (bool, error) {
		keepTrying := iterations[i].keepTrying
		err := iterations[i].err
		i++
		return keepTrying, err
	}
}

func Test_backoff_bad_MaxSleepTime(t *testing.T) {
	require.Panics(t, func() {
		_ = Backoff(echo(true, nil), BackoffOptions{
			MaxSleepTime: 0,
		})
	})

}

func Test_backoff_bad_InitialGapSize(t *testing.T) {
	require.Panics(t, func() {
		_ = Backoff(echo(true, nil), BackoffOptions{
			MaxSleepTime:   1 * time.Second,
			InitialGapSize: 0,
		})
	})
}

func Test_backoff_bad_MaxJitterSize(t *testing.T) {
	require.Panics(t, func() {
		_ = Backoff(echo(true, nil), BackoffOptions{
			MaxSleepTime:   1 * time.Second,
			InitialGapSize: 1 * time.Millisecond,
			MaxJitterSize:  -1 * time.Millisecond,
		})
	})
}

func Test_backoff_bad_MaxJitterSize_large(t *testing.T) {
	require.Panics(t, func() {
		_ = Backoff(echo(true, nil), BackoffOptions{
			MaxSleepTime:   1 * time.Second,
			InitialGapSize: 1 * time.Millisecond,
			MaxJitterSize:  501 * time.Millisecond,
		})
	})
}

func Test_Backoff_first_success(t *testing.T) {
	sleeper := libtime.NewSleeperMock(t)
	defer sleeper.MinimockFinish()

	opts := BackoffOptions{
		MaxSleepTime:   32 * time.Millisecond,
		InitialGapSize: 1 * time.Millisecond,
		Sleeper:        sleeper,
	}

	err := Backoff(echo(true, nil), opts)
	require.NoError(t, err)
}

func Test_Backoff_later_success(t *testing.T) {
	sleeper := libtime.NewSleeperMock(t)
	defer sleeper.MinimockFinish()

	invocation := 0
	sleeper.SleepMock.Set(func(d time.Duration) {
		invocation++
		switch invocation {
		case 1:
			require.Equal(t, 1*time.Millisecond, d)
		case 2:
			require.Equal(t, 2*time.Millisecond, d)
		case 3:
			require.Equal(t, 4*time.Millisecond, d)
		default:
			t.Fatalf("unexpected invocation: %d", invocation)
		}
	})

	f := schedule([]iteration{
		{true, errOops},
		{true, errOops},
		{true, errOops},
		{true, nil},
	})

	opts := BackoffOptions{
		MaxSleepTime:   32 * time.Millisecond,
		InitialGapSize: 1 * time.Millisecond,
		Sleeper:        sleeper,
	}

	err := Backoff(f, opts)
	require.NoError(t, err)
}

func Test_Backoff_giveup_partway(t *testing.T) {
	sleeper := libtime.NewSleeperMock(t)
	defer sleeper.MinimockFinish()

	invocation := 0
	sleeper.SleepMock.Set(func(d time.Duration) {
		invocation++
		switch invocation {
		case 1:
			require.Equal(t, 1*time.Millisecond, d)
		case 2:
			require.Equal(t, 2*time.Millisecond, d)
		case 3:
			require.Equal(t, 4*time.Millisecond, d)
		case 4:
			require.Equal(t, 8*time.Millisecond, d)
		default:
			t.Fatalf("unexpected invocation: %d", invocation)
		}
	})

	f := schedule([]iteration{
		{true, errOops},
		{true, errOops},
		{true, errOops},
		{true, errOops},
		{false, errOops},
	})

	opts := BackoffOptions{
		MaxSleepTime:   32 * time.Millisecond,
		InitialGapSize: 1 * time.Millisecond,
		Sleeper:        sleeper,
	}

	err := Backoff(f, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "instructed to stop retrying")
}

func Test_Backoff_all_fail(t *testing.T) {
	sleeper := libtime.NewSleeperMock(t)
	defer sleeper.MinimockFinish()

	invocation := 0
	sleeper.SleepMock.Set(func(d time.Duration) {
		invocation++
		switch invocation {
		case 1:
			require.Equal(t, 1*time.Millisecond, d)
		case 2:
			require.Equal(t, 2*time.Millisecond, d)
		case 3:
			require.Equal(t, 4*time.Millisecond, d)
		case 4:
			require.Equal(t, 8*time.Millisecond, d)
		case 5:
			require.Equal(t, 16*time.Millisecond, d)
		case 6:
			// 6th iteration gets truncated
			require.Equal(t, 1*time.Millisecond, d)
		default:
			t.Fatalf("unexpected invocation: %d", invocation)
		}
	})

	opts := BackoffOptions{
		MaxSleepTime:   32 * time.Millisecond,
		InitialGapSize: 1 * time.Millisecond,
		Sleeper:        sleeper,
	}

	err := Backoff(echo(true, errOops), opts)
	require.Equal(t, ErrMaximumTimeExceeded, err)
}

func Test_Backoff_jitter(t *testing.T) {
	sleeper := libtime.NewSleeperMock(t)
	defer sleeper.MinimockFinish()

	invocation := 0
	sleeper.SleepMock.Set(func(d time.Duration) {
		invocation++
		switch invocation {
		case 1:
			require.Equal(t, time.Duration(3354690), d)
		case 2:
			require.Equal(t, time.Duration(4043278), d)
		case 3:
			require.Equal(t, time.Duration(5042819), d)
		case 4:
			require.Equal(t, time.Duration(8985913), d)
		case 5:
			require.Equal(t, time.Duration(10573300), d)
		default:
			t.Fatalf("unexpected invocation: %d", invocation)
		}
	})

	opts := BackoffOptions{
		MaxSleepTime:   32 * time.Millisecond,
		InitialGapSize: 1 * time.Millisecond,
		MaxJitterSize:  5 * time.Millisecond,
		RandomSeed:     666, // deterministic randomness
		Sleeper:        sleeper,
	}

	err := Backoff(echo(true, errOops), opts)
	require.Equal(t, ErrMaximumTimeExceeded, err)
}

func Test_Backoff_real_life(t *testing.T) {
	opts := BackoffOptions{
		MaxSleepTime:   300 * time.Millisecond,
		InitialGapSize: 10 * time.Millisecond,
		MaxJitterSize:  10 * time.Millisecond,
		RandomSeed:     time.Now().Unix(),
		// Sleeper will default to time.Sleep
	}

	f := schedule([]iteration{
		{true, errOops}, // 10 + [0-10] (lower: 10, upper: 20)
		{true, errOops}, // 20 + [0-10] (lower: 30, upper: 50)
		{true, errOops}, // 40 + [0-10] (lower: 70, upper: 100)
		{true, errOops}, // 80 + [0-10] (lower: 150, upper: 190)
		{true, nil},
	})

	start := time.Now()
	err := Backoff(f, opts)
	elapsed := time.Since(start)

	require.NoError(t, err)

	// lower bound
	require.True(t, elapsed >= 150*time.Millisecond)

	// upper bound
	require.True(t, elapsed <= 190*time.Millisecond)
}
