package blind_injection

import (
	"github.com/SergeyAkaikin/vulnscan/internal/sqli"
)

type Blind struct {
	sqli.Base
}

func New(url, method string, parameters ...string) *Blind {
	return &Blind{
		sqli.Base{
			Url:         url,
			Method:      method,
			QParameters: parameters,
		},
	}
}
