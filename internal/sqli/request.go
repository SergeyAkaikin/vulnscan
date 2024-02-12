package sqli

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Base struct {
	Url         string
	Method      string
	QParameters []string
}

func (b *Base) MakeRequest(parameterInd int, query string, cookie *http.Cookie) (res *http.Response, payload string, err error) {
	switch b.Method {
	case http.MethodGet:
		payload = b.uriQuery(parameterInd, query)
		req, err := http.NewRequest(http.MethodGet, payload, nil)
		if err != nil {
			return nil, "", err
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}
		res, err = http.DefaultClient.Do(req)

	case http.MethodPost:
		payload = b.bodyQuery(parameterInd, query)
		buff := bytes.NewReader([]byte(payload))
		req, err := http.NewRequest(http.MethodPost, b.Url, buff)
		if err != nil {
			return nil, "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if cookie != nil {
			req.AddCookie(cookie)
		}
		res, err = http.DefaultClient.Do(req)
	}
	return
}

func (b *Base) uriQuery(ind int, query string) string {
	reqStr := strings.Builder{}
	reqStr.WriteString(fmt.Sprintf("%s?%s%s", b.Url, b.QParameters[ind], query))
	for i := 0; i < len(b.QParameters); i++ {
		if i == ind {
			continue
		}

		reqStr.WriteString(fmt.Sprintf("&%s", b.QParameters[i]))
	}
	return reqStr.String()
}

func (b *Base) bodyQuery(ind int, query string) string {
	bodyStr := strings.Builder{}
	bodyStr.WriteString(fmt.Sprintf("%s%s", b.QParameters[ind], query))
	for i := 0; i < len(b.QParameters); i++ {
		if i == ind {
			continue
		}
		bodyStr.WriteString(fmt.Sprintf("&%s", b.QParameters[i]))
	}

	return bodyStr.String()
}
