package error_injection

import (
	"github.com/SergeyAkaikin/vulnscan/internal/sqli"
	"io"
	"net/http"
	"regexp"
)

type ErrorBased struct {
	sqli.Base
}

func New(url, method string, parameters ...string) *ErrorBased {
	return &ErrorBased{
		sqli.Base{
			Url:         url,
			Method:      method,
			QParameters: parameters,
		},
	}
}

func (e *ErrorBased) Inject(extCookies *http.Cookie) (injectable bool, payload string) {
	for i := 0; i < len(e.QParameters); i++ {
		res, payload, err := e.MakeRequest(i, "'", extCookies)
		if err != nil {
			continue
		}
		content, err := io.ReadAll(res.Body)
		if err != nil {
			continue
		}

		reg := regexp.MustCompile(sqli.ErrorPayload)
		if reg.Match(content) {
			return true, payload
		}
	}

	return false, ""
}
