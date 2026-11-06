package sample

import "testing"

// TestIntentionalRed exists only to prove that CI fails on a red test.
// It is removed again in the next commit (P1.T3.S3.M2).
func TestIntentionalRed(t *testing.T) {
	t.Fatal("intentional failure: CI must catch this")
}
