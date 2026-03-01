package wait

import (
	"context"
	"errors"
	"time"
)

// CheckFunc is a function that checks if a condition has been met.
// It returns:
//   - done: true if the condition is met and polling should stop
//   - err: an error if the check failed (this will stop polling and return the error)
type CheckFunc func() (done bool, err error)

// WaitForCondition polls a condition function at the specified interval until either:
//   - The condition returns done=true (success)
//   - The condition returns an error (failure)
//   - The context deadline is exceeded (timeout)
//
// The interval parameter specifies how long to wait between checks.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
//	defer cancel()
//
//	err := wait.WaitForCondition(ctx, 10*time.Second, func() (bool, error) {
//	    resp, err := client.GetServerByUuid(ctx, uuid)
//	    if err != nil {
//	        return false, err
//	    }
//	    return resp.Settings.IsReachable && resp.Settings.IsUsable, nil
//	})
func WaitForCondition(ctx context.Context, interval time.Duration, checkFn CheckFunc) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Check immediately before first tick
	done, err := checkFn()
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			done, err := checkFn()
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}
}

// TimeoutError creates a user-friendly error message for context deadline exceeded errors
func TimeoutError(operation string, timeout time.Duration) error {
	return errors.New(operation + " did not complete within " + timeout.String())
}
