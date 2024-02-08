package blind_injection

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Blind struct {
	Url         string
	Method      string
	QParameters []string
}

func New(url, method string, parameters ...string) *Blind {
	return &Blind{
		Url:         url,
		Method:      method,
		QParameters: parameters,
	}
}

func (b *Blind) makeRequest(parId int, query string) (res *http.Response, payload string, err error) {
	switch b.Method {
	case http.MethodGet:
		payload = b.uriQuery(parId, query)
		res, err = http.Get(payload)
	case http.MethodPost:
		payload = b.bodyQuery(parId, query)
		buff := bytes.NewReader([]byte(payload))
		res, err = http.Post(b.Url, "application/x-www-form-urlencoded", buff)
	}
	return
}

func (b *Blind) uriQuery(ind int, query string) string {
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

func (b *Blind) bodyQuery(ind int, query string) string {
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
