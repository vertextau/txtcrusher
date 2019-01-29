// Package pastebin provides methods for working with Pastebin service.
package pastebin

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
)

const (
	baseURLLogin  = "https://pastebin.com/api/api_login.php"
	baseURLPastes = "https://pastebin.com/api/api_raw.php"
	baseURLPost   = "https://pastebin.com/api/api_post.php"
	baseURLRaw    = "https://pastebin.com/raw"
)

// Modificators for a paste.
const (
	Public = iota
	Unlisted
	Private
)

// API synonyms.
const (
	apiDevKey       = "api_dev_key"
	apiOption       = "api_option"
	apiPasteCode    = "api_paste_code"
	apiPasteExpDate = "api_paste_expire_date"
	apiPasteFormat  = "api_paste_format"
	apiPasteKey     = "api_paste_key"
	apiPasteName    = "api_paste_name"
	apiPastePrivate = "api_paste_private"
	apiResultLimit  = "api_results_limit"
	apiUserKey      = "api_user_key"
	apiUserName     = "api_user_name"
	apiUserPassword = "api_user_password"
)

// The expiration date of a paste.
const (
	ExpNever     = "N"
	Exp10Minutes = "10M"
	Exp1Hour     = "1H"
	Exp1Day      = "1D"
	Exp1Week     = "1W"
	Exp2Weeks    = "2W"
	Exp1Month    = "1M"
	Exp6Months   = "6M"
	Exp1Year     = "1Y"
)

// Command options.
const (
	optionDeletePaste = "delete"
	optionList        = "list"
	optionPaste       = "paste"
	optionShowPaste   = "show_paste"
	optionUserInfo    = "userdetails"
)

// Regexp for bad API response.
const (
	badResponseRegexp = "^Bad API request"
)

var (
	badRespRegexp = regexp.MustCompile(badResponseRegexp)
)

// Main structure for storing developer and user keys.
type (
	Pastebin struct {
		DeveloperKey string
		UserKey      string
	}
)

// Handles requests to a server.
func RequestHandler(baseURL string, v *url.Values) (*string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := string(respBody)
	if checkResponseError(&result) {
		return nil, errors.New(result)
	}

	return &result, nil
}

// Checks for a bad response from a server.
func checkResponseError(resp *string) bool {
	if badRespRegexp.MatchString(*resp) {
		return true
	} else {
		return false
	}
}

// Creates an 'api_user_key' using the api member login system.
func (p Pastebin) GetUserKey(username, password string) (*string, error) {
	v := url.Values{}

	v.Set(apiDevKey, p.DeveloperKey)
	v.Set(apiUserName, username)
	v.Set(apiUserPassword, password)

	result, err := RequestHandler(baseURLLogin, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Creates a new paste.
func (p Pastebin) CreateNewPaste(pasteText *string, guest bool, pasteTitle, pasteFormat, expireDate string, pasteMod int) (*string, error) {
	v := url.Values{}

	v.Set(apiDevKey, p.DeveloperKey)

	if !guest {
		v.Set(apiUserKey, p.UserKey)
	}

	v.Set(apiPasteCode, *pasteText)
	v.Set(apiPasteName, pasteTitle)
	v.Set(apiPastePrivate, strconv.Itoa(pasteMod))
	v.Set(apiPasteExpDate, expireDate)
	v.Set(apiPasteFormat, pasteFormat)
	v.Set(apiOption, optionPaste)

	result, err := RequestHandler(baseURLPost, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Gets a users information and settings.
func (p Pastebin) GetUserInfo() (*string, error) {
	v := url.Values{}

	v.Set(apiDevKey, p.DeveloperKey)
	v.Set(apiUserKey, p.UserKey)
	v.Set(apiOption, optionUserInfo)

	result, err := RequestHandler(baseURLPost, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Lists pastes created by a user.
func (p Pastebin) ListUserPastes(resultLimit int) (*string, error) {
	if resultLimit == 0 {
		resultLimit = 50
	} else if resultLimit < 0 || resultLimit > 1000 {
		return nil, errors.New("Results limit must be in range [1, 1000]")
	}

	v := url.Values{}

	v.Set(apiDevKey, p.DeveloperKey)
	v.Set(apiUserKey, p.UserKey)
	v.Set(apiResultLimit, strconv.Itoa(resultLimit))
	v.Set(apiOption, optionList)

	result, err := RequestHandler(baseURLPost, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Gets raw paste output of users pastes including 'private' pastes.
func (p Pastebin) GetUserPaste(pasteKey string) (*string, error) {
	v := url.Values{}
	v.Set(apiDevKey, p.DeveloperKey)
	v.Set(apiUserKey, p.UserKey)
	v.Set(apiPasteKey, pasteKey)
	v.Set(apiOption, optionShowPaste)

	result, err := RequestHandler(baseURLPastes, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Removes a paste created by a user.
func (p Pastebin) DeleteUserPaste(pasteKey string) (*string, error) {
	v := url.Values{}
	v.Set(apiDevKey, p.DeveloperKey)
	v.Set(apiUserKey, p.UserKey)
	v.Set(apiPasteKey, pasteKey)
	v.Set(apiOption, optionDeletePaste)

	result, err := RequestHandler(baseURLPost, &v)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Gets raw paste output of any 'public' and 'unlisted' pastes.
func GetPaste(pasteKey string) (*string, error) {
	client := &http.Client{}

	baseURL, err := url.Parse(baseURLRaw)
	if err != nil {
		return nil, err
	}

	baseURL.Path = path.Join(baseURL.Path, pasteKey)

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	paste := string(respBody)

	return &paste, nil
}
