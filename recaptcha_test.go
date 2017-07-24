// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness

package recaptcha

import (
	"errors"
	"testing"
)

func Test_checkReturn(t *testing.T) {
	tests := []struct {
		resp       interface{}
		err        error
		shouldPass bool
	}{
		{
			resp:       &recaptchaResponse{},
			err:        nil,
			shouldPass: true,
		},
		{
			resp:       nil,
			err:        nil,
			shouldPass: false,
		},
		{
			resp:       nil,
			err:        errors.New("hey"),
			shouldPass: false,
		},
	}
	for k, tt := range tests {
		err := checkReturn("test", tt.resp, tt.err)
		if !tt.shouldPass && err == nil {
			t.Errorf("test #%d should fail", k)
		} else if tt.shouldPass {
			if err != nil {
				t.Errorf("test #%d should pass", k)
			}
		}
	}
}
