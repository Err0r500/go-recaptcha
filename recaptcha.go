// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness

package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RecaptchaBridge struct {
	KeyGetter
	address string
}

type KeyGetter interface {
	getSecretKey() (string, error)
}

type recaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []int     `json:"error-codes"`
}

func NewRecaptchaBridge(kG KeyGetter, address ...string) *RecaptchaBridge {
	rB := RecaptchaBridge{KeyGetter: kG}
	if len(address) == 0 {
		rB.address = "https://www.google.com/recaptcha/api/siteverify"
	}
	if len(address) > 1 {
		return nil
	}
	return &rB
}

// Confirm is the public facade.
// Pass it the response and optionnally the client IP
/*
USAGE :
rB := NewRecaptchaBridge(simpleGetter{})
ok, err := Confirm(response, remoteip)
if err != nil {
	-> the validation process encountered a technical issue
}
if ok {
	-> the response to the callenge is fine
}
*/
func (rB *RecaptchaBridge) Confirm(response string, remoteip ...string) (ok bool, err error) {
	query, err := rB.setupQuery(response, remoteip)
	if err := checkReturn("rB.setupQuery()", query, err); err != nil {
		return false, err
	}

	attempt, err := rB.check(query)
	if err = checkReturn("rB.check()", attempt, err); err != nil {
		return false, err
	}
	return attempt.Success, nil
}

func (rB *RecaptchaBridge) check(req *url.Values) (r *recaptchaResponse, err error) {
	resp, err := http.PostForm(rB.address, *req)
	if err = checkReturn("http.PostForm()", resp, err); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r = &recaptchaResponse{Success: false}
	if err := r.fromJSON(body); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *recaptchaResponse) fromJSON(body []byte) error {
	return json.Unmarshal(body, r)
}

func checkReturn(funcName string, resp interface{}, err error) error {
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("%s returned a nil pointer", funcName)
	}
	return nil
}
