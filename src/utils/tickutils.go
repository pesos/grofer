package utils

import (
	"context"
	"time"
)

func TickUntilDone(ctx context.Context, refreshRate int64, action func() error) error {
	ticker := time.NewTicker(time.Duration(refreshRate) * time.Millisecond)
	defer ticker.Stop()

	for {
		err := action()
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err() // Stop execution if end signal received
		case <-ticker.C:
			break
		}
	}
}
