package recaptcha

import (
	"errors"
	"os"
)

type simpleKeyGetter struct {
	secret string
}

func (sKG simpleKeyGetter) getSecretKey() (string, error) {
	if os.Getenv("CAPTCHA_SECRET") != "" {
		return os.Getenv("CAPTCHA_SECRET"), nil
	}
	if sKG.secret != "" {
		return sKG.secret, nil
	}
	return "", errors.New("can't fetch the secret key")
}
