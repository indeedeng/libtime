package libtime

import (
	"context"
	"testing"
	"time"
)

func Test_After(t *testing.T) {
	<-After(context.Background(), 3*time.Nanosecond)
}

func Test_After_canceledCtx(t *testing.T) {
	const afterTimeout = 1 * time.Millisecond

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	after := After(canceledCtx, afterTimeout)

	afterX2 := time.NewTimer(2 * afterTimeout)
	defer afterX2.Stop()

	// This test is racy. If it ever fails, consider increasing the multiplier or the timeout, or both
	select {
	case <-after:
		t.Fatal("no messages should be sent on the channel returned by After, if the context is canceled")
	case <-afterX2.C:
	}
}
