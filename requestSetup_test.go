package recaptcha

import (
	"net/url"
	"testing"
)

func TestSetupQuery(t *testing.T) {
	tests := []struct {
		req        recaptchaRequest
		expected   url.Values
		shouldPass bool
	}{
		{
			req:        recaptchaRequest{Response: "foo", RemoteIP: "122.2.2.2"},
			expected:   map[string][]string{"secret": []string{"testSecretKey"}, "response": []string{"foo"}, "remoteip": []string{"122.2.2.2"}},
			shouldPass: true,
		},
		{
			req:        recaptchaRequest{Response: "foo"},
			expected:   map[string][]string{"secret": []string{"testSecretKey"}, "response": []string{"foo"}},
			shouldPass: true,
		},
		{
			req:        recaptchaRequest{},
			expected:   map[string][]string{},
			shouldPass: false,
		},
	}
	rB := NewRecaptchaBridge(simpleKeyGetter{secret: "testSecretKey"})
	for k, tt := range tests {
		returned, err := rB.setupQuery(tt.req.Response, []string{tt.req.RemoteIP})
		if !tt.shouldPass && err == nil {
			t.Errorf("test #%d should fail", k)
		} else if tt.shouldPass {
			if err != nil {
				t.Errorf("test #%d should pass", k)
			}
			if len(*returned) != len(tt.expected) {
				t.Fatalf("wrong return on test #%d :\n%s\n%s", k, tt.expected, *returned)
			}
			for i, v := range tt.expected {
				if v[0] != (*returned)[i][0] {
					t.Errorf("wrong return on test #%d :\n%s\n%s", k, tt.expected, *returned)
				}
			}
		}
	}

	_, err := rB.setupQuery(tests[0].req.Response, []string{})
	if err != nil {
		t.Errorf("setupQuery should pass without remoteip")
	}
}
