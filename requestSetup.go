package recaptcha

import (
	"errors"
	"net/url"

	"github.com/google/go-querystring/query"
)

type recaptchaRequest struct {
	SecretKey string `url:"secret"`
	RemoteIP  string `url:"remoteip,omitempty"`
	Response  string `url:"response"`
}

func (rB RecaptchaBridge) setupQuery(response string, remoteip []string) (*url.Values, error) {
	secretKey, err := rB.KeyGetter.getSecretKey()
	if err != nil {
		return nil, err
	}

	remoteIP := "" // will be ommited in the request if not set
	if len(remoteip) == 1 {
		remoteIP = remoteip[0]
	}
	req := recaptchaRequest{
		RemoteIP:  remoteIP,
		Response:  response,
		SecretKey: secretKey,
	}

	if req.Response == "" || req.SecretKey == "" {
		return nil, errors.New("a mandatory property is missing")
	}
	uV, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	return &uV, nil
}
