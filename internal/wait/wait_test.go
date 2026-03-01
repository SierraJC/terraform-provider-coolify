package wait

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitForCondition_ImmediateSuccess(t *testing.T) {
	ctx := context.Background()
	
	callCount := 0
	checkFn := func() (bool, error) {
		callCount++
		return true, nil
	}
	
	err := WaitForCondition(ctx, 100*time.Millisecond, checkFn)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if callCount != 1 {
		t.Errorf("expected checkFn to be called once, got %d calls", callCount)
	}
}

func TestWaitForCondition_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()
	
	callCount := 0
	checkFn := func() (bool, error) {
		callCount++
		// Succeed on the third call
		return callCount >= 3, nil
	}
	
	err := WaitForCondition(ctx, 10*time.Millisecond, checkFn)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if callCount < 3 {
		t.Errorf("expected at least 3 calls, got %d", callCount)
	}
}

func TestWaitForCondition_ErrorFromCheck(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("check failed")
	
	checkFn := func() (bool, error) {
		return false, expectedErr
	}
	
	err := WaitForCondition(ctx, 100*time.Millisecond, checkFn)
	
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestWaitForCondition_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	checkFn := func() (bool, error) {
		// Never succeeds
		return false, nil
	}
	
	err := WaitForCondition(ctx, 10*time.Millisecond, checkFn)
	
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestWaitForCondition_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	callCount := 0
	checkFn := func() (bool, error) {
		callCount++
		if callCount == 2 {
			cancel() // Cancel after second call
		}
		return false, nil
	}
	
	err := WaitForCondition(ctx, 10*time.Millisecond, checkFn)
	
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestWaitForCondition_ErrorOnSecondCall(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("transient error")
	
	callCount := 0
	checkFn := func() (bool, error) {
		callCount++
		if callCount == 2 {
			return false, expectedErr
		}
		return false, nil
	}
	
	err := WaitForCondition(ctx, 10*time.Millisecond, checkFn)
	
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	
	// Should have called exactly twice before erroring
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestTimeoutError(t *testing.T) {
	tests := []struct {
		operation string
		timeout   time.Duration
		expected  string
	}{
		{
			operation: "server validation",
			timeout:   3 * time.Minute,
			expected:  "server validation did not complete within 3m0s",
		},
		{
			operation: "service deployment",
			timeout:   10 * time.Minute,
			expected:  "service deployment did not complete within 10m0s",
		},
		{
			operation: "custom operation",
			timeout:   30 * time.Second,
			expected:  "custom operation did not complete within 30s",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			err := TimeoutError(tt.operation, tt.timeout)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}
		})
	}
}

func TestWaitForCondition_Timing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timing test in short mode")
	}
	
	ctx := context.Background()
	interval := 50 * time.Millisecond
	expectedCalls := 5
	
	callCount := 0
	start := time.Now()
	
	checkFn := func() (bool, error) {
		callCount++
		return callCount >= expectedCalls, nil
	}
	
	err := WaitForCondition(ctx, interval, checkFn)
	elapsed := time.Since(start)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if callCount != expectedCalls {
		t.Errorf("expected %d calls, got %d", expectedCalls, callCount)
	}
	
	// Should take at least (expectedCalls-1) * interval
	// -1 because first call is immediate
	minDuration := time.Duration(expectedCalls-1) * interval
	if elapsed < minDuration {
		t.Errorf("expected at least %v, took %v", minDuration, elapsed)
	}
	
	// Should not take more than 2x the expected duration (generous margin)
	maxDuration := minDuration * 2
	if elapsed > maxDuration {
		t.Errorf("expected no more than %v, took %v", maxDuration, elapsed)
	}
}
