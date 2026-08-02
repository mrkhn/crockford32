package crockford32

import (
	"fmt"
	"strings"
	"testing"
)

func TestSecurityVulnerabilityReproduction(t *testing.T) {
	t.Run("ErrInputTooLong_TruncationAndQuoting", func(t *testing.T) {
		longInput := strings.Repeat("A", 1024)
		_, err := Parse(longInput)
		if err == nil {
			t.Fatal("expected error for long input, got nil")
		}

		errStr := err.Error()
		// New behavior: should be quoted and truncated.
		// "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA..." (32 As)
		expectedInput := strings.Repeat("A", 32) + "..."
		expectedInfix := fmt.Sprintf(`Parse(%q):`, expectedInput)
		if !strings.Contains(errStr, expectedInfix) {
			t.Errorf("expected error message to contain truncated and quoted input, but it didn't. Got: %s", errStr)
		}
	})

	t.Run("ErrInvalidRunes_SafeQuoting", func(t *testing.T) {
		maliciousInput := "invalid\ninput\""
		_, err := Parse(maliciousInput)
		if err == nil {
			t.Fatal("expected error for invalid input, got nil")
		}

		errStr := err.Error()
		// New behavior: should be quoted. %q will escape \n and ".
		expectedInput := maliciousInput
		expectedInfix := fmt.Sprintf(`Parse(%q):`, expectedInput)
		if !strings.Contains(errStr, expectedInfix) {
			t.Errorf("expected error message to contain safely quoted input, but it didn't. Got: %s", errStr)
		}
	})
}
