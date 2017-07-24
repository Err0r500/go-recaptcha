package recaptcha

import (
	"os"
	"testing"
)

func TestGetSecretKey(t *testing.T) {
	tests := []struct {
		fetchMode  string
		expected   string
		shouldPass bool
	}{
		{
			fetchMode:  "env",
			expected:   "fooSecretKey",
			shouldPass: true,
		},
		{
			expected:   "fooSecretKey",
			shouldPass: true,
		},
		{
			fetchMode:  "env",
			expected:   "",
			shouldPass: false,
		},
	}
	for k, tt := range tests {
		if err := os.Unsetenv("CAPTCHA_SECRET"); err != nil {
			t.Fatal("failed to unset CAPTCHA_SECRET")
		}
		sKG := simpleKeyGetter{}
		if tt.fetchMode == "env" {
			if err := os.Setenv("CAPTCHA_SECRET", tt.expected); err != nil {
				t.Fatal("failed to set CAPTCHA_SECRET")
			}
		} else {
			sKG.secret = tt.expected
		}
		returned, err := sKG.getSecretKey()
		if !tt.shouldPass && err == nil {
			t.Errorf("test #%d should fail", k)
		} else if tt.shouldPass {
			if err != nil {
				t.Errorf("test #%d should pass", k)
			}
			if returned != tt.expected {
				t.Errorf("wrong return on test #%d :\n%s\n%s", k, tt.expected, returned)
			}
		}
	}
}
